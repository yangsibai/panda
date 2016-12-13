package helper

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime"
	"path/filepath"
	"time"
)

func UploadObjectToOss2(key string, filepath string, extension string) error {
	client, err := oss.New(Config.Oss2Endpoint, Config.Oss2AccessKeyId, Config.Oss2AccessKeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(Config.Oss2BucketName)
	if err != nil {
		return err
	}

	return bucket.PutObjectFromFile(key, filepath, oss.ContentType(mime.TypeByExtension(extension)))
}

func GetOss2SavedKey(uuid string) string {
	date := time.Now().Format("20060102")
	return filepath.Join("image", date, uuid)
}
