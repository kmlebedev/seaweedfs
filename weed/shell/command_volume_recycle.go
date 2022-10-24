package shell

import (
	"context"
	"flag"
	"fmt"
	"github.com/seaweedfs/seaweedfs/weed/operation"
	"github.com/seaweedfs/seaweedfs/weed/pb"
	"github.com/seaweedfs/seaweedfs/weed/pb/volume_server_pb"
	"github.com/seaweedfs/seaweedfs/weed/storage/needle"
	"io"
	"sort"
)

func init() {
	Commands = append(Commands, &commandVolumeRecycle{})
}

type commandVolumeRecycle struct {
	writer       io.Writer
	applyRecycle *bool
}

func (c *commandVolumeRecycle) Name() string {
	return "volume.recycle"
}

func (c *commandVolumeRecycle) Help() string {
	return `volume.recycle -freeThreshold=0.3
	This command commandVolumeRecycle, when the cluster free storage more than ${freeThreshold}, 
     it will trigger the deletion of the oldest ${recycleVolumeCounter} file
`
}

func (c *commandVolumeRecycle) Do(args []string, commandEnv *CommandEnv, writer io.Writer) (err error) {
	recycleCommand := flag.NewFlagSet(c.Name(), flag.ContinueOnError)
	freeThreshold := recycleCommand.Float64("freeThreshold", 0.3, "recycle when free is more than this limit")
	c.applyRecycle = recycleCommand.Bool("force", false, "apply to recycle volumes")
	if err = recycleCommand.Parse(args); err != nil {
		return nil
	}
	if err = commandEnv.confirmIsLocked(args); err != nil {
		return
	}
	infoAboutSimulationMode(writer, *c.applyRecycle, "-force")

	c.writer = writer
	topologyInfo, _, err := collectTopologyInfo(commandEnv, 0)
	if err != nil {
		return
	}
	dataCenterInfo := topologyInfo.DataCenterInfos
	volumeIdToVolumeMap := make(map[uint32]pb.ServerAddress)
	var volumeIds []uint32
	var volumeServers []pb.ServerAddress

	for _, dataCenter := range dataCenterInfo {
		if dataCenter.RackInfos == nil || len(dataCenter.RackInfos) == 0 {
			fmt.Fprintf(c.writer, "Error dataCenter rack is empty\n")
			continue
		}
		for _, rack := range dataCenter.RackInfos {
			if rack.DataNodeInfos == nil || len(rack.DataNodeInfos) == 0 {
				fmt.Fprintf(c.writer, "Error BuildClusterVo DataNodeInfos == nil || len(vr.DataNodeInfos) == 0\n")
				continue
			}
			for _, dataNode := range rack.DataNodeInfos {
				volumeServers = append(volumeServers, pb.NewServerAddressFromDataNode(dataNode))
				for _, disk := range dataNode.DiskInfos {
					if disk.VolumeInfos == nil || len(disk.VolumeInfos) == 0 {
						fmt.Fprintf(c.writer, "Error disk.VolumeInfos == nil || len(disk.VolumeInfos) == 0\n")
						continue
					}
					for _, volume := range disk.VolumeInfos {
						volumeIdToVolumeMap[volume.Id] = pb.NewServerAddressFromDataNode(dataNode)
						volumeIds = append(volumeIds, volume.Id)
					}
				}
			}
		}
	}

	sort.Slice(volumeIds, func(i, j int) bool {
		if volumeIds[i] < volumeIds[j] {
			return true
		}
		return false
	})
	diskStatus, errorDiskStatus := volumeDisk(volumeServers, commandEnv)
	if errorDiskStatus != nil {
		fmt.Fprintf(c.writer, "Error %+v\n", errorDiskStatus.Error())
		return errorDiskStatus
	}
	freePer := float64(diskStatus.Free) / float64(diskStatus.All)
	fmt.Fprintf(c.writer, "Free:%d, all:%d, freePer:%f\n", diskStatus.Free, diskStatus.All, freePer)
	if freePer >= *freeThreshold {
		for _, volumeId := range volumeIds {
			volumeServer := volumeIdToVolumeMap[volumeId]
			err := deleteVolume(commandEnv.option.GrpcDialOption, needle.VolumeId(volumeId), volumeServer)
			if err != nil {
				fmt.Fprintf(c.writer, "Error deleteVolume %+v volumeId is %d  %s\n", volumeServer, volumeId, err.Error())
				return err
			}
			fmt.Fprintf(c.writer, "deleteVolume %+v volumeId is %d success\n", volumeServer, volumeId)

		}
	}
	fmt.Fprintf(c.writer, "VolumeRecycle do success\n")
	return nil
}

/*
Get cluster storage information
*/
func volumeDisk(volumeServers []pb.ServerAddress, commandEnv *CommandEnv) (diskStatus volume_server_pb.DiskStatus, err error) {

	var diskAll uint64
	var diskFree uint64
	var diskUsed uint64
	for _, volumeServer := range volumeServers {
		err := operation.WithVolumeServerClient(false, volumeServer, commandEnv.option.GrpcDialOption, func(volumeServerClient volume_server_pb.VolumeServerClient) error {
			resp, statusErr := volumeServerClient.VolumeServerStatus(context.Background(), &volume_server_pb.VolumeServerStatusRequest{})
			if statusErr != nil {
				return statusErr
			}
			if resp.DiskStatuses == nil || len(resp.DiskStatuses) == 0 {
				return fmt.Errorf("%+v Disk is empty", volumeServer)
			}
			for _, disk := range resp.DiskStatuses {
				diskFree += disk.Free
				diskAll += disk.All
				diskUsed += disk.Used
			}
			return nil
		})
		if err != nil {
			return diskStatus, err
		}
	}
	diskStatus.All = diskAll
	diskStatus.Used = diskUsed
	diskStatus.Free = diskFree
	return diskStatus, nil
}
