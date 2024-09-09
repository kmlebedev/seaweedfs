package weed_server

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/seaweedfs/seaweedfs/weed/topology"

	"github.com/seaweedfs/raft"

	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/seaweedfs/seaweedfs/weed/pb/master_pb"
	"github.com/seaweedfs/seaweedfs/weed/security"
	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
	"github.com/seaweedfs/seaweedfs/weed/storage/super_block"
	"github.com/seaweedfs/seaweedfs/weed/storage/types"
)

func (ms *MasterServer) DoAutomaticVolumeGrow(req *topology.VolumeGrowRequest) {
	glog.V(1).Infoln("starting automatic volume grow")
	start := time.Now()
	newVidLocations, err := ms.vg.AutomaticGrowByType(req.Option, ms.grpcDialOption, ms.Topo, req.Count)
	glog.V(1).Infoln("finished automatic volume grow, cost ", time.Now().Sub(start))
	if err != nil {
		glog.V(1).Infof("automatic volume grow failed: %+v", err)
		return
	}
	for _, newVidLocation := range newVidLocations {
		ms.broadcastToClients(&master_pb.KeepConnectedResponse{VolumeLocation: newVidLocation})
	}
}

func (ms *MasterServer) ProcessGrowRequest() {
	go func() {
		for {
			if !ms.Topo.IsLeader() {
				continue
			}
			dcs := ms.Topo.ListDataCenters()
			for _, vlc := range ms.Topo.ListVolumeLayoutCollections() {
				vl := vlc.VolumeLayout
				if vl.HasGrowRequest() {
					continue
				}
				if vl.ShouldGrowVolumes(vlc.Collection) {
					vl.AddGrowRequest()
					ms.volumeGrowthRequestChan <- &topology.VolumeGrowRequest{
						Option: vlc.ToGrowOption(),
						Count:  vl.GetLastGrowCount(),
						Reason: "collection autogrow",
					}
				} else {
					for _, dc := range dcs {
						if vl.ShouldGrowVolumesByDataNode("DataCenter", dc) {
							vl.AddGrowRequest()
							volumeGrowOption := vlc.ToGrowOption()
							volumeGrowOption.DataCenter = dc
							ms.volumeGrowthRequestChan <- &topology.VolumeGrowRequest{
								Option: volumeGrowOption,
								Count:  vl.GetLastGrowCount(),
								Force:  true,
								Reason: "per-dc autogrow",
							}
						}
					}
				}
			}
			time.Sleep(14*time.Minute + time.Duration(120*rand.Float32())*time.Second)
		}
	}()
	go func() {
		filter := sync.Map{}
		for {
			req, ok := <-ms.volumeGrowthRequestChan
			if !ok {
				break
			}

			option := req.Option
			vl := ms.Topo.GetVolumeLayout(option.Collection, option.ReplicaPlacement, option.Ttl, option.DiskType)

			if !ms.Topo.IsLeader() {
				//discard buffered requests
				time.Sleep(time.Second * 1)
				vl.DoneGrowRequest()
				continue
			}

			// filter out identical requests being processed
			found := false
			filter.Range(func(k, v interface{}) bool {
				existingReq := k.(*topology.VolumeGrowRequest)
				if existingReq.Equals(req) {
					found = true
				}
				return !found
			})

			// not atomic but it's okay
			if found || (!req.Force && !vl.ShouldGrowVolumes(req.Option.Collection)) {
				glog.V(4).Infoln("discard volume grow request")
				time.Sleep(time.Millisecond * 211)
				vl.DoneGrowRequest()
				continue
			}

			filter.Store(req, nil)
			// we have lock called inside vg
			glog.V(0).Infof("volume grow %+v", req)
			go func(req *topology.VolumeGrowRequest, vl *topology.VolumeLayout) {
				ms.DoAutomaticVolumeGrow(req)
				vl.DoneGrowRequest()
				filter.Delete(req)
			}(req, vl)
		}
	}()
}

func (ms *MasterServer) LookupVolume(ctx context.Context, req *master_pb.LookupVolumeRequest) (*master_pb.LookupVolumeResponse, error) {

	resp := &master_pb.LookupVolumeResponse{}
	volumeLocations := ms.lookupVolumeId(req.VolumeOrFileIds, req.Collection)

	for _, volumeOrFileId := range req.VolumeOrFileIds {
		vid := volumeOrFileId
		commaSep := strings.Index(vid, ",")
		if commaSep > 0 {
			vid = vid[0:commaSep]
		}
		if result, found := volumeLocations[vid]; found {
			var locations []*master_pb.Location
			for _, loc := range result.Locations {
				locations = append(locations, &master_pb.Location{
					Url:        loc.Url,
					PublicUrl:  loc.PublicUrl,
					DataCenter: loc.DataCenter,
					GrpcPort:   uint32(loc.GrpcPort),
				})
			}
			var auth string
			if commaSep > 0 { // this is a file id
				auth = string(security.GenJwtForVolumeServer(ms.guard.SigningKey, ms.guard.ExpiresAfterSec, result.VolumeOrFileId))
			}
			resp.VolumeIdLocations = append(resp.VolumeIdLocations, &master_pb.LookupVolumeResponse_VolumeIdLocation{
				VolumeOrFileId: result.VolumeOrFileId,
				Locations:      locations,
				Error:          result.Error,
				Auth:           auth,
			})
		}
	}

	return resp, nil
}

func (ms *MasterServer) Statistics(ctx context.Context, req *master_pb.StatisticsRequest) (*master_pb.StatisticsResponse, error) {

	if !ms.Topo.IsLeader() {
		return nil, raft.NotLeaderError
	}

	if req.Replication == "" {
		req.Replication = ms.option.DefaultReplicaPlacement
	}
	replicaPlacement, err := super_block.NewReplicaPlacementFromString(req.Replication)
	if err != nil {
		return nil, err
	}
	ttl, err := needle.ReadTTL(req.Ttl)
	if err != nil {
		return nil, err
	}

	volumeLayout := ms.Topo.GetVolumeLayout(req.Collection, replicaPlacement, ttl, types.ToDiskType(req.DiskType))
	stats := volumeLayout.Stats()
	totalSize := ms.Topo.GetDiskUsages().GetMaxVolumeCount() * int64(ms.option.VolumeSizeLimitMB) * 1024 * 1024
	resp := &master_pb.StatisticsResponse{
		TotalSize: uint64(totalSize),
		UsedSize:  stats.UsedSize,
		FileCount: stats.FileCount,
	}

	return resp, nil
}

func (ms *MasterServer) VolumeList(ctx context.Context, req *master_pb.VolumeListRequest) (*master_pb.VolumeListResponse, error) {

	if !ms.Topo.IsLeader() {
		return nil, raft.NotLeaderError
	}

	resp := &master_pb.VolumeListResponse{
		TopologyInfo:      ms.Topo.ToTopologyInfo(),
		VolumeSizeLimitMb: uint64(ms.option.VolumeSizeLimitMB),
	}

	return resp, nil
}

func (ms *MasterServer) LookupEcVolume(ctx context.Context, req *master_pb.LookupEcVolumeRequest) (*master_pb.LookupEcVolumeResponse, error) {

	if !ms.Topo.IsLeader() {
		return nil, raft.NotLeaderError
	}

	resp := &master_pb.LookupEcVolumeResponse{}

	ecLocations, found := ms.Topo.LookupEcShards(needle.VolumeId(req.VolumeId))

	if !found {
		return resp, fmt.Errorf("ec volume %d not found", req.VolumeId)
	}

	resp.VolumeId = req.VolumeId

	for shardId, shardLocations := range ecLocations.Locations {
		var locations []*master_pb.Location
		for _, dn := range shardLocations {
			locations = append(locations, &master_pb.Location{
				Url:        string(dn.Id()),
				PublicUrl:  dn.PublicUrl,
				DataCenter: dn.GetDataCenterId(),
			})
		}
		resp.ShardIdLocations = append(resp.ShardIdLocations, &master_pb.LookupEcVolumeResponse_EcShardIdLocation{
			ShardId:   uint32(shardId),
			Locations: locations,
		})
	}

	return resp, nil
}

func (ms *MasterServer) VacuumVolume(ctx context.Context, req *master_pb.VacuumVolumeRequest) (*master_pb.VacuumVolumeResponse, error) {

	if !ms.Topo.IsLeader() {
		return nil, raft.NotLeaderError
	}

	resp := &master_pb.VacuumVolumeResponse{}

	ms.Topo.Vacuum(ms.grpcDialOption, float64(req.GarbageThreshold), ms.option.MaxParallelVacuumPerServer, req.VolumeId, req.Collection, ms.preallocateSize)

	return resp, nil
}

func (ms *MasterServer) DisableVacuum(ctx context.Context, req *master_pb.DisableVacuumRequest) (*master_pb.DisableVacuumResponse, error) {

	ms.Topo.DisableVacuum()
	resp := &master_pb.DisableVacuumResponse{}
	return resp, nil
}

func (ms *MasterServer) EnableVacuum(ctx context.Context, req *master_pb.EnableVacuumRequest) (*master_pb.EnableVacuumResponse, error) {

	ms.Topo.EnableVacuum()
	resp := &master_pb.EnableVacuumResponse{}
	return resp, nil
}

func (ms *MasterServer) VolumeMarkReadonly(ctx context.Context, req *master_pb.VolumeMarkReadonlyRequest) (*master_pb.VolumeMarkReadonlyResponse, error) {

	if !ms.Topo.IsLeader() {
		return nil, raft.NotLeaderError
	}

	resp := &master_pb.VolumeMarkReadonlyResponse{}

	replicaPlacement, _ := super_block.NewReplicaPlacementFromByte(byte(req.ReplicaPlacement))
	vl := ms.Topo.GetVolumeLayout(req.Collection, replicaPlacement, needle.LoadTTLFromUint32(req.Ttl), types.ToDiskType(req.DiskType))
	dataNodes := ms.Topo.Lookup(req.Collection, needle.VolumeId(req.VolumeId))

	for _, dn := range dataNodes {
		if dn.Ip == req.Ip && dn.Port == int(req.Port) {
			if req.IsReadonly {
				vl.SetVolumeReadOnly(dn, needle.VolumeId(req.VolumeId))
			} else {
				vl.SetVolumeWritable(dn, needle.VolumeId(req.VolumeId))
			}
		}
	}

	return resp, nil
}

func (ms *MasterServer) VolumeGrow(ctx context.Context, req *master_pb.AssignRequest) (*master_pb.VolumeGrowResponse, error) {
	if !ms.Topo.IsLeader() {
		return nil, raft.NotLeaderError
	}
	if req.Replication == "" {
		req.Replication = ms.option.DefaultReplicaPlacement
	}
	replicaPlacement, err := super_block.NewReplicaPlacementFromString(req.Replication)
	if err != nil {
		return nil, err
	}
	ttl, err := needle.ReadTTL(req.Ttl)
	if err != nil {
		return nil, err
	}
	volumeGrowOption := topology.VolumeGrowOption{
		Collection:         req.Collection,
		ReplicaPlacement:   replicaPlacement,
		Ttl:                ttl,
		DiskType:           types.ToDiskType(req.DiskType),
		Preallocate:        ms.preallocateSize,
		DataCenter:         req.DataCenter,
		Rack:               req.Rack,
		DataNode:           req.DataNode,
		MemoryMapMaxSizeMb: req.MemoryMapMaxSizeMb,
	}
	volumeGrowRequest := topology.VolumeGrowRequest{
		Option: &volumeGrowOption,
		Count:  req.WritableVolumeCount,
		Force:  true,
		Reason: "grpc volume grow",
	}
	replicaCount := int64(req.WritableVolumeCount * uint32(replicaPlacement.GetCopyCount()))

	if ms.Topo.AvailableSpaceFor(&volumeGrowOption) < replicaCount {
		return nil, fmt.Errorf("only %d volumes left, not enough for %d", ms.Topo.AvailableSpaceFor(&volumeGrowOption), replicaCount)
	}

	if !ms.Topo.DataCenterExists(volumeGrowOption.DataCenter) {
		err = fmt.Errorf("data center %v not found in topology", volumeGrowOption.DataCenter)
	}

	ms.DoAutomaticVolumeGrow(&volumeGrowRequest)

	return &master_pb.VolumeGrowResponse{}, nil
}
