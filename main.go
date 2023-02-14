package main

import (
	"flag"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
	"path/filepath"
)

const (
	Output       = "build"
	MatchPattern = "assets/*/*"

	//ENDPOINT     = "https://oss-cn-hangzhou.aliyuncs.com"
	//BUCKET       = "frontend-weekly"
	//AccessKeyId     = "LTAI5tMFnXMeAycS1pSBfsKw"
	//AccessKeySecret = "cpJQIwUfCtdM61hyIu2r4tsP70pSV5"
)

var (
	endpoint = flag.String("endpoint", "", "填写Bucket对应的Endpoint")
	bucket   = flag.String("bucket", "", "填写存储空间名称")

	AccessKeyId     = os.Getenv("ACCESS_KEY_ID")
	AccessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
)

func main() {
	flag.Parse()

	// 创建OSSClient实例
	client, err := oss.New(*endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		log.Fatal(err)
	}
	// 获取存储空间
	bucket, err := client.Bucket(*bucket)
	if err != nil {
		log.Fatal(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	buildPath := filepath.Join(pwd, Output)
	matches, err := filepath.Glob(filepath.Join(buildPath, MatchPattern))
	if err != nil {
		log.Fatal(err)
	}

	totalCount, ignoreCount, uploadCount := len(matches), 0, 0

	for _, match := range matches {
		// 计算相对路径作为 objectKey
		objectKey, err := filepath.Rel(buildPath, match)
		if err != nil {
			log.Fatal(err)
		}
		isExist, err := bucket.IsObjectExist(objectKey)
		if err != nil {
			log.Fatal(err)
		}
		if isExist {
			ignoreCount++
			continue
		}
		// 上传文件
		if err = bucket.PutObjectFromFile(objectKey, match); err != nil {
			log.Fatal(err)
		}
		uploadCount++
	}

	log.Printf(
		"[success] Total: %d Uploaded: %d Ignored: %d",
		totalCount,
		uploadCount,
		ignoreCount,
	)
}
