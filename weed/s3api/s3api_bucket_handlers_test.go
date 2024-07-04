package s3api

import (
	"github.com/seaweedfs/seaweedfs/weed/s3api/s3err"
	"testing"
	"time"
)

func TestListBucketsHandler(t *testing.T) {

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult><Owner><ID></ID><DisplayName></DisplayName></Owner><Buckets><Bucket><Name>test1</Name><CreationDate>2011-04-09 12:34:49 +0000 UTC</CreationDate></Bucket><Bucket><Name>test2</Name><CreationDate>2011-02-09 12:34:49 +0000 UTC</CreationDate></Bucket></Buckets></ListAllMyBucketsResult>`
	var response ListAllMyBucketsResult

	var bucketsList ListAllMyBucketsList
	bucketsList.Bucket = append(bucketsList.Bucket, &ListAllMyBucketsEntry{
		Name:         "test1",
		CreationDate: time.Date(2011, 4, 9, 12, 34, 49, 0, time.UTC).String(),
	})
	bucketsList.Bucket = append(bucketsList.Bucket, &ListAllMyBucketsEntry{
		Name:         "test2",
		CreationDate: time.Date(2011, 2, 9, 12, 34, 49, 0, time.UTC).String(),
	})

	response = ListAllMyBucketsResult{
		Owner: &CanonicalUser{
			ID:          "",
			DisplayName: "",
		},
		Buckets: &bucketsList,
	}

	encoded := string(s3err.EncodeXMLResponse(response))
	if encoded != expected {
		t.Errorf("unexpected output:%s\nexpecting:%s", encoded, expected)
	}
}
