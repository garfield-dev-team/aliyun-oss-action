package main

import (
	"flag"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
	"path/filepath"
)

// TODO: 除了 ENDPOINT，秘钥信息通过 GitHub 账号配置
const (
	Output       = "build"
	MatchPattern = "assets/*/*"
)

var (
	endpoint = flag.String("endpoint", "", "Path to main Go main package.")
	bucket   = flag.String("bucket", "", "Override action name, the default name is the package name.")

	AccessKeyId     = os.Getenv("ACCESS_KEY_ID")
	AccessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
)

func main() {
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
			log.Println(objectKey, "资源已存在，忽略上传")
			continue
		}
		// 上传文件
		if err = bucket.PutObjectFromFile(objectKey, match); err != nil {
			log.Fatal(err)
		}
		log.Println(objectKey, "资源上传成功")
	}
}
