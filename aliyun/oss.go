package aliyun

import (
	"fmt"
	"io"
	"os"
	path2 "path"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSS struct {
	EndPoint   string
	BucketName string
	Aliyun
}

func NewOSS(aliyun Aliyun, EndPoint, Bucket string) *OSS {
	oss := &OSS{
		EndPoint:   EndPoint,
		BucketName: Bucket,
		Aliyun:     aliyun,
	}
	return oss
}
func (o OSS) Upload(path string, reader io.Reader) (string, error) {
	client, err := oss.New(o.EndPoint, o.AK, o.AS)
	if err != nil {
		return "", err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(o.BucketName)
	if err != nil {
		return "", err
	}
	storageType := oss.ObjectStorageClass(oss.StorageStandard)

	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	err = bucket.PutObject(path, reader, storageType, objectAcl)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("http://%s.%s/%s", o.BucketName, o.EndPoint, path)
	return url, nil
}
func (o OSS) UploadFile(dir, filepath string) (string, error) {
	fileIO, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer fileIO.Close()
	path := path2.Join(dir, path2.Base(filepath))
	return o.Upload(path, fileIO)
}
