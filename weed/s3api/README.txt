see https://blog.aqwari.net/xml-schema-go/

1. go install github.com/xuri/xgen/cmd/xgen@latest
2. Add EncodingType element for ListBucketResult in AmazonS3.xsd
3. xgen -i AmazonS3.xsd -o s3api_xsd_generated.go -l Go -p s3api
4. Remove empty Grantee struct in s3api_xsd_generated.go
