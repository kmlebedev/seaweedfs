package s3api

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/seaweedfs/seaweedfs/weed/filer"
	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/seaweedfs/seaweedfs/weed/pb"
	"github.com/seaweedfs/seaweedfs/weed/pb/filer_pb"
	"github.com/seaweedfs/seaweedfs/weed/pb/s3_pb"
	"github.com/seaweedfs/seaweedfs/weed/s3api/s3_constants"
	"github.com/seaweedfs/seaweedfs/weed/s3api/s3err"
	"github.com/seaweedfs/seaweedfs/weed/util"
	"net/http"
	"path/filepath"
	"sync"
)

type BucketAccessControlPolicies struct {
	sync.RWMutex
	policies *map[string]*BucketAccessControlPolicy
}

type BucketAccessControlPolicy struct {
	sync.RWMutex
	s3.AccessControlPolicy

	ObjectOwnership string
	Owner           *s3.Owner `type:"structure"` // Container for the bucket owner's display name and ID.
	// A list of grants for access controls.
	Grants []*s3.Grant `locationName:"AccessControlList" locationNameList:"Grant" type:"list"`
}

func NewBucketAccessControlPolicies() *BucketAccessControlPolicies {
	policies := make(map[string]*BucketAccessControlPolicy)
	bacp := BucketAccessControlPolicies{
		policies: &policies,
	}
	return &bacp
}

func (bacp *BucketAccessControlPolicies) LoadConfigurationFromBytes(content *[]byte) error {
	bpCfg := &s3_pb.S3BucketAccessControlPolices{}
	if err := filer.ParseS3ConfigurationFromBytes(content, bpCfg); err != nil {
		glog.Warningf("unmarshal bucket policies error: %v", err)
		return fmt.Errorf("unmarshal bucket policies error: %v", err)
	}
	bacp.loadPbConfig(bpCfg)
	return nil
}

func (s3a *S3ApiServer) SaveBucketAccessControlPoliciesConfig() (err error) {
	buf := bytes.Buffer{}
	cfg := s3a.bacp.ToPbConfig()
	if err := filer.ProtoToText(&buf, cfg); err != nil {
		return fmt.Errorf("ProtoToText: %s", err)
	}
	return pb.WithGrpcFilerClient(false, 0, s3a.option.Filer, s3a.option.GrpcDialOption, func(client filer_pb.SeaweedFilerClient) error {
		err = util.Retry("saveACLs", func() error {
			return filer.SaveInsideFiler(client, s3_constants.S3ConfigDir, s3_constants.BucketACLsConfigFile, buf.Bytes())
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func (bacp *BucketAccessControlPolicies) ToPbConfig() *s3_pb.S3BucketAccessControlPolices {
	cfg := s3_pb.S3BucketAccessControlPolices{
		Policies: make(map[string]*s3_pb.AccessControlPolicy),
	}
	bacp.RLock()
	defer bacp.RUnlock()
	if bacp.policies == nil {
		return &cfg
	}
	for bucket, acp := range *bacp.policies {
		acp.RLock()
		var grants []*s3_pb.Grant
		for _, grant := range acp.Grants {
			if grant == nil {
				continue
			}
			pbGrant := s3_pb.Grant{}
			if grant.Permission != nil {
				pbGrant.Permission = *grant.Permission
			}
			if grant.Grantee != nil {
				pbGrant.Grantee = &s3_pb.Grantee{}
				if grant.Grantee.DisplayName != nil {
					pbGrant.Grantee.DisplayName = *grant.Grantee.DisplayName
				}
				if grant.Grantee.EmailAddress != nil {
					pbGrant.Grantee.EmailAddress = *grant.Grantee.EmailAddress
				}
				if grant.Grantee.ID != nil {
					pbGrant.Grantee.Id = *grant.Grantee.ID
				}
				if grant.Grantee.Type != nil {
					pbGrant.Grantee.Type = *grant.Grantee.Type
				}
				if grant.Grantee.URI != nil {
					pbGrant.Grantee.Uri = *grant.Grantee.URI
				}
			}
			grants = append(grants, &pbGrant)
		}
		pbAcp := s3_pb.AccessControlPolicy{
			ObjectOwnership:   acp.ObjectOwnership,
			AccessControlList: grants,
		}
		if acp.Owner != nil {
			pbAcp.Owner = &s3_pb.BucketOwner{}
			if acp.Owner.DisplayName != nil {
				pbAcp.Owner.DisplayName = *acp.Owner.DisplayName
			}
			if acp.Owner.ID != nil {
				pbAcp.Owner.Id = *acp.Owner.ID
			}
		}
		cfg.Policies[bucket] = &pbAcp
		acp.RUnlock()
	}
	return &cfg
}

func (bacp *BucketAccessControlPolicies) loadPbConfig(cfg *s3_pb.S3BucketAccessControlPolices) {
	bucketPolicies := make(map[string]*BucketAccessControlPolicy)
	for bucket, policy := range cfg.Policies {
		var grants []*s3.Grant
		for _, grant := range policy.AccessControlList {
			grants = append(grants, GrantFromPb(grant))
		}
		bucketPolicies[bucket] = &BucketAccessControlPolicy{
			ObjectOwnership: policy.ObjectOwnership,
			Owner: &s3.Owner{
				DisplayName: aws.String(policy.Owner.DisplayName),
				ID:          aws.String(policy.Owner.Id),
			},
			Grants: grants,
		}
	}
	bacp.Lock()
	bacp.policies = &bucketPolicies
	bacp.Unlock()
}

func (bacp *BucketAccessControlPolicies) DeleteAccessControlPolicy(bucket string) {
	bacp.Lock()
	defer bacp.Unlock()
	if _, ok := (*bacp.policies)[bucket]; ok {
		delete(*bacp.policies, bucket)
	}
}

func (bacp *BucketAccessControlPolicies) GetAccessControlPolicy(bucket string) *BucketAccessControlPolicy {
	bacp.RLock()
	acp, ok := (*bacp.policies)[bucket]
	bacp.RUnlock()
	if !ok {
		acp = &BucketAccessControlPolicy{
			Owner:  &s3.Owner{},
			Grants: []*s3.Grant{},
		}
		bacp.Lock()
		(*bacp.policies)[bucket] = acp
		bacp.Unlock()
	}
	return acp
}

func (bacp *BucketAccessControlPolicies) GetOwnerAccountId(bucket string) string {
	if ownerAccountId := bacp.GetAccessControlPolicy(bucket).Owner.ID; ownerAccountId != nil {
		return *ownerAccountId
	}
	return ""
}

func getAccountId(r *http.Request) string {
	id := r.Header.Get(s3_constants.AmzAccountId)
	if len(id) == 0 {
		return s3account.AccountAnonymous.Id
	} else {
		return id
	}
}

func (s3a *S3ApiServer) checkAccessByOwnership(r *http.Request, acp *BucketAccessControlPolicy) s3err.ErrorCode {
	requestAccountId := getAccountId(r)
	if acp != nil && ValidateAccount(requestAccountId, *acp.Owner.ID) {
		return s3err.ErrNone
	}
	return s3err.ErrAccessDenied
}

// Todo Check access for PutBucketAclHandler
func (s3a *S3ApiServer) checkAccessForPutBucketAcl(requestAccountId, bucket string) (*BucketMetaData, s3err.ErrorCode) {
	bucketMetadata, errCode := s3a.bucketRegistry.GetBucketMetadata(bucket)
	if errCode != s3err.ErrNone {
		return nil, errCode
	}

	if bucketMetadata.ObjectOwnership == s3_constants.OwnershipBucketOwnerEnforced {
		return nil, s3err.AccessControlListNotSupported
	}

	if ValidateAccount(requestAccountId, *bucketMetadata.Owner.ID) {
		return bucketMetadata, s3err.ErrNone
	}

	if len(bucketMetadata.Acl) > 0 {
		reqGrants := DetermineRequiredGrants(requestAccountId, s3_constants.PermissionWriteAcp)
		for _, bucketGrant := range bucketMetadata.Acl {
			for _, reqGrant := range reqGrants {
				if GrantEquals(bucketGrant, reqGrant) {
					return bucketMetadata, s3err.ErrNone
				}
			}
		}
	}
	glog.V(3).Infof("acl denied! request account id: %s", requestAccountId)
	return nil, s3err.ErrAccessDenied
}

func updateBucketEntry(s3a *S3ApiServer, entry *filer_pb.Entry) error {
	return s3a.updateEntry(s3a.option.BucketsPath, entry)
}

// Check Bucket/BucketAcl Read related access
// includes:
// - HeadBucketHandler
// - GetBucketAclHandler
// - ListObjectsV1Handler
// - ListObjectsV2Handler
// - ListMultipartUploadsHandler
func (s3a *S3ApiServer) checkAccessForReadBucket(r *http.Request, bucket, aclAction string) (*BucketMetaData, s3err.ErrorCode) {
	bucketMetadata, errCode := s3a.bucketRegistry.GetBucketMetadata(bucket)
	if errCode != s3err.ErrNone {
		return nil, errCode
	}

	if bucketMetadata.ObjectOwnership == s3_constants.OwnershipBucketOwnerEnforced {
		return bucketMetadata, s3err.ErrNone
	}

	requestAccountId := GetAccountId(r)
	if ValidateAccount(requestAccountId, *bucketMetadata.Owner.ID) {
		return bucketMetadata, s3err.ErrNone
	}

	if len(bucketMetadata.Acl) > 0 {
		reqGrants := DetermineRequiredGrants(requestAccountId, aclAction)
		for _, bucketGrant := range bucketMetadata.Acl {
			for _, reqGrant := range reqGrants {
				if GrantEquals(bucketGrant, reqGrant) {
					return bucketMetadata, s3err.ErrNone
				}
			}
		}
	}

	glog.V(3).Infof("acl denied! request account id: %s", requestAccountId)
	return nil, s3err.ErrAccessDenied
}

// Todo Check ObjectAcl-Read related access
// includes:
// - GetObjectAclHandler
func (s3a *S3ApiServer) checkAccessForReadObjectAcl(r *http.Request, bucket, object string) (acp *s3.AccessControlPolicy, errCode s3err.ErrorCode) {
	bucketMetadata, errCode := s3a.bucketRegistry.GetBucketMetadata(bucket)
	if errCode != s3err.ErrNone {
		return nil, errCode
	}

	getAcpFunc := func() (*s3.AccessControlPolicy, s3err.ErrorCode) {
		entry, err := getObjectEntry(s3a, bucket, object)
		if err != nil {
			if err == filer_pb.ErrNotFound {
				return nil, s3err.ErrNoSuchKey
			} else {
				return nil, s3err.ErrInternalError
			}
		}
		if entry.IsDirectory {
			return nil, s3err.ErrExistingObjectIsDirectory
		}
		acpOwnerId := GetAcpOwner(entry.Extended, *bucketMetadata.Owner.ID)
		acpOwnerName := s3a.accountManager.IdNameMapping[acpOwnerId]
		acpGrants := GetAcpGrants(&acpOwnerId, entry.Extended)
		acp = &s3.AccessControlPolicy{
			Owner: &s3.Owner{
				ID:          &acpOwnerId,
				DisplayName: &acpOwnerName,
			},
			Grants: acpGrants,
		}
		return acp, s3err.ErrNone
	}
	if bucketMetadata.ObjectOwnership == s3_constants.OwnershipBucketOwnerEnforced {
		return getAcpFunc()
	}
	requestAccountId := GetAccountId(r)
	acp, errCode = getAcpFunc()
	if errCode != s3err.ErrNone {
		return nil, errCode
	}
	if ValidateAccount(requestAccountId, *acp.Owner.ID) {
		return acp, s3err.ErrNone
	}
	if acp.Grants != nil {
		reqGrants := DetermineRequiredGrants(requestAccountId, s3_constants.PermissionReadAcp)
		for _, requiredGrant := range reqGrants {
			for _, grant := range acp.Grants {
				if GrantEquals(requiredGrant, grant) {
					return acp, s3err.ErrNone
				}
			}
		}
	}
	glog.V(3).Infof("CheckAccessForReadObjectAcl denied! request account id: %s", requestAccountId)
	return nil, s3err.ErrAccessDenied
}

// Check Object-Read related access
// includes:
// - GetObjectHandler
//
// offload object access validation to Filer layer
// - s3acl.CheckObjectAccessForReadObject
func (s3a *S3ApiServer) checkBucketAccessForReadObject(r *http.Request, bucket string) s3err.ErrorCode {
	bucketMetadata, errCode := s3a.bucketRegistry.GetBucketMetadata(bucket)
	if errCode != s3err.ErrNone {
		return errCode
	}

	if bucketMetadata.ObjectOwnership != s3_constants.OwnershipBucketOwnerEnforced {
		//offload object acl validation to filer layer
		_, defaultErrorCode := s3a.checkAccessForReadBucket(r, bucket, s3_constants.PermissionRead)
		if defaultErrorCode != s3err.ErrNone {
			r.Header.Set(s3_constants.XSeaweedFSHeaderAmzBucketAccessDenied, "true")
		}
		r.Header.Set(s3_constants.XSeaweedFSHeaderAmzBucketOwnerId, *bucketMetadata.Owner.ID)
	}

	return s3err.ErrNone
}

// Todo Check ObjectAcl-Write related access
// includes:
// - PutObjectAclHandler
func (s3a *S3ApiServer) checkAccessForWriteObjectAcl(r *http.Request, bucket, object string) (*filer_pb.Entry, string, []*s3.Grant, s3err.ErrorCode) {
	//bucketMetadata, errCode := s3a.bucketRegistry.GetBucketMetadata(bucket)
	//if errCode != s3err.ErrNone {
	//	return nil, "", nil, errCode
	//}

	requestAccountId := GetAccountId(r)
	//acp := s3a.bacp.GetAccessControlPolicy(bucket)
	//reqOwnerId, grants, errCode := s3acl.ExtractObjectAcl(r, s3a.accountManager, acp.ObjectOwnership, *acp.Owner.ID, requestAccountId, false)
	//if errCode != s3err.ErrNone {
	//	return nil, "", nil, errCode
	//}

	//if bucketMetadata.ObjectOwnership == s3_constants.OwnershipBucketOwnerEnforced {
	//	return nil, "", nil, s3err.AccessControlListNotSupported
	//}

	//object acl
	// objectEntry, err := getObjectEntry(s3a, bucket, object)
	// if err != nil {
	//	if err == filer_pb.ErrNotFound {
	//		return nil, "", nil, s3err.ErrNoSuchKey
	//	}
	//	return nil, "", nil, s3err.ErrInternalError
	//}

	//if objectEntry.IsDirectory {
	//	return nil, "", nil, s3err.ErrExistingObjectIsDirectory
	//}

	//objectOwner := s3acl.GetAcpOwner(objectEntry.Extended, *bucketMetadata.Owner.ID)
	//object owner is immutable

	// Todo use s3test
	//if reqOwnerId != "" && reqOwnerId != *acp.Owner.ID {
	//	return nil, "", nil, s3err.ErrAccessDenied
	//}
	//if s3acl.ValidateAccount(requestAccountId, objectOwner) {
	//	return objectEntry, objectOwner, grants, s3err.ErrNone
	//}

	//objectGrants := s3acl.GetAcpGrants(nil, objectEntry.Extended)
	//if objectGrants != nil {
	//	requiredGrants := s3acl.DetermineRequiredGrants(requestAccountId, s3_constants.PermissionWriteAcp)
	//	for _, objectGrant := range objectGrants {
	//		for _, requiredGrant := range requiredGrants {
	//			if s3acl.GrantEquals(objectGrant, requiredGrant) {
	//				return objectEntry, objectOwner, grants, s3err.ErrNone
	//			}
	//		}
	//	}
	//}

	glog.V(3).Infof("checkAccessForWriteObjectAcl denied! request account id: %s", requestAccountId)
	return nil, "", nil, s3err.ErrAccessDenied
}

func updateObjectEntry(s3a *S3ApiServer, bucket, object string, entry *filer_pb.Entry) error {
	dir, _ := filepath.Split(object)
	return s3a.updateEntry(util.Join(s3a.option.BucketsPath, bucket, dir), entry)
}

// CheckAccessForPutObject Check ACL for PutObject API
// includes:
// - PutObjectHandler
func (s3a *S3ApiServer) CheckAccessForPutObject(r *http.Request, bucket, object string) s3err.ErrorCode {
	return s3a.checkAccessForWriteObject(r, bucket, object, GetAccountId(r))
}

// CheckAccessForPutObjectPartHandler Check Acl for Upload object part
// includes:
// - PutObjectPartHandler
func (s3a *S3ApiServer) CheckAccessForPutObjectPartHandler(r *http.Request, bucket string) s3err.ErrorCode {
	bucketMetadata, errCode := s3a.bucketRegistry.GetBucketMetadata(bucket)
	if errCode != s3err.ErrNone {
		return errCode
	}
	if bucketMetadata.ObjectOwnership == s3_constants.OwnershipBucketOwnerEnforced {
		return s3err.ErrNone
	}
	accountId := GetAccountId(r)
	if !CheckBucketAccess(accountId, bucketMetadata, s3_constants.PermissionWrite) {
		return s3err.ErrAccessDenied
	}
	return s3err.ErrNone
}

// CheckAccessForNewMultipartUpload Check Acl for  API
// includes:
// - NewMultipartUploadHandler
func (s3a *S3ApiServer) CheckAccessForNewMultipartUpload(r *http.Request, bucket, object string) (s3err.ErrorCode, string) {
	accountId := GetAccountId(r)
	if accountId == IdentityAnonymous.AccountId {
		return s3err.ErrAccessDenied, ""
	}
	errCode := s3a.checkAccessForWriteObject(r, bucket, object, accountId)
	return errCode, accountId
}

func (s3a *S3ApiServer) CheckAccessForAbortMultipartUpload(r *http.Request, bucket, object string) s3err.ErrorCode {
	return s3a.CheckAccessWithBucketOwnerAndInitiator(r, bucket, object)
}

func (s3a *S3ApiServer) CheckAccessForCompleteMultipartUpload(r *http.Request, bucket, object string) s3err.ErrorCode {
	bucketMetadata, errCode := s3a.bucketRegistry.GetBucketMetadata(bucket)
	if errCode != s3err.ErrNone {
		return errCode
	}

	if bucketMetadata.ObjectOwnership != s3_constants.OwnershipBucketOwnerEnforced {
		accountId := getAccountId(r)
		if !CheckBucketAccess(accountId, bucketMetadata, s3_constants.PermissionWrite) {
			return s3err.ErrAccessDenied
		}
	}
	return s3err.ErrNone
}

func (s3a *S3ApiServer) CheckAccessForListMultipartUploadParts(r *http.Request, bucket, object string) s3err.ErrorCode {
	return s3a.CheckAccessWithBucketOwnerAndInitiator(r, bucket, object)
}

// CheckAccessWithBucketOwnerAndInitiator Check Access Permission with 'bucketOwner' and 'multipartUpload initiator'
func (s3a *S3ApiServer) CheckAccessWithBucketOwnerAndInitiator(r *http.Request, bucket, object string) s3err.ErrorCode {
	bucketMetadata, errCode := s3a.bucketRegistry.GetBucketMetadata(bucket)
	if errCode != s3err.ErrNone {
		return errCode
	}

	//bucket access allowed
	accountId := GetAccountId(r)
	if ValidateAccount(*bucketMetadata.Owner.ID, accountId) {
		return s3err.ErrNone
	}

	//multipart initiator allowed
	entry, err := getMultipartUpload(s3a, bucket, object)
	if err != nil {
		if err != filer_pb.ErrNotFound {
			return s3err.ErrInternalError
		}
	} else {
		uploadInitiator, ok := entry.Extended[s3_constants.ExtAmzMultipartInitiator]
		if !ok || accountId == string(uploadInitiator) {
			return s3err.ErrNone
		}
	}
	glog.V(3).Infof("CheckAccessWithBucketOwnerAndInitiator denied! request account id: %s", accountId)
	return s3err.ErrAccessDenied
}

func (s3a *S3ApiServer) checkAccessForWriteObject(r *http.Request, bucket, object, requestAccountId string) s3err.ErrorCode {
	bucketMetadata, errCode := s3a.bucketRegistry.GetBucketMetadata(bucket)
	if errCode != s3err.ErrNone {
		return errCode
	}
	requestOwnerId, grants, errCode := ExtractObjectAcl(r, s3a.accountManager, bucketMetadata.ObjectOwnership, *bucketMetadata.Owner.ID, requestAccountId, true)
	if errCode != s3err.ErrNone {
		return errCode
	}

	if bucketMetadata.ObjectOwnership == s3_constants.OwnershipBucketOwnerEnforced {
		return s3err.ErrNone
	}

	if !CheckBucketAccess(requestAccountId, bucketMetadata, s3_constants.PermissionWrite) {
		return s3err.ErrAccessDenied
	}

	if requestOwnerId == "" {
		requestOwnerId = requestAccountId
	}
	entry, err := getObjectEntry(s3a, bucket, object)
	if err != nil {
		if err == filer_pb.ErrNotFound {
			SetAcpOwnerHeader(r, requestOwnerId)
			SetAcpGrantsHeader(r, grants)
			return s3err.ErrNone
		}
		return s3err.ErrInternalError
	}

	objectOwnerId := GetAcpOwner(entry.Extended, *bucketMetadata.Owner.ID)
	//object owner is immutable
	if requestOwnerId != "" && objectOwnerId != requestOwnerId {
		return s3err.ErrAccessDenied
	}

	//Only the owner of the bucket and the owner of the object can overwrite the object
	if ValidateAccount(requestOwnerId, objectOwnerId, *bucketMetadata.Owner.ID) {
		glog.V(3).Infof("checkAccessForWriteObject denied! request account id: %s, expect account id: %s", requestAccountId, *bucketMetadata.Owner.ID)
		return s3err.ErrAccessDenied
	}

	SetAcpOwnerHeader(r, objectOwnerId)
	SetAcpGrantsHeader(r, grants)
	return s3err.ErrNone
}

func CheckBucketAccess(requestAccountId string, bucketMetadata *BucketMetaData, permission string) bool {
	if ValidateAccount(requestAccountId, *bucketMetadata.Owner.ID) {
		return true
	} else {
		if len(bucketMetadata.Acl) > 0 {
			reqGrants := DetermineRequiredGrants(requestAccountId, permission)
			for _, bucketGrant := range bucketMetadata.Acl {
				for _, requiredGrant := range reqGrants {
					if GrantEquals(bucketGrant, requiredGrant) {
						return true
					}
				}
			}
		}
	}
	glog.V(3).Infof("CheckBucketAccess denied! request account id: %s", requestAccountId)
	return false
}

func getObjectEntry(s3a *S3ApiServer, bucket, object string) (*filer_pb.Entry, error) {
	return s3a.getEntry(util.Join(s3a.option.BucketsPath, bucket), object)
}

func getMultipartUpload(s3a *S3ApiServer, bucket, object string) (*filer_pb.Entry, error) {
	return s3a.getEntry(s3a.genUploadsFolder(bucket), object)
}
