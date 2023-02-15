package main

import (
	"flag"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sourcegraph/conc/pool"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"
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

func init() {
	// 解决 log 时区问题
	var cstZone = time.FixedZone("CST", 8*3600) // 东八区
	time.Local = cstZone
}

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

	totalCount := len(matches)
	var ignoreCount, uploadCount int64
	start := time.Now()
	log.Println("[info] uploading...")

	log.Println("===", runtime.NumCPU())

	p := pool.New().
		WithMaxGoroutines(runtime.NumCPU() * 2)
	for _, match := range matches {
		match := match
		p.Go(func() {
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
				atomic.AddInt64(&ignoreCount, 1)
				return
			}
			// 上传文件
			if err = bucket.PutObjectFromFile(objectKey, match); err != nil {
				log.Fatal(err)
			}
			atomic.AddInt64(&uploadCount, 1)
		})
	}

	p.Wait()

	log.Printf(
		"[success] total: %d uploaded: %d ignored: %d in %.2fs",
		totalCount,
		uploadCount,
		ignoreCount,
		time.Since(start).Seconds(),
	)
}
