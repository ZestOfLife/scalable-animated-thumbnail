package main

import (
	"CompileWorker/queue"
	"CompileWorker/senders"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/minio/minio-go/v7"
)

// https://stackoverflow.com/questions/43073681/golang-binary-search
func binarySearch(a []string, search int) (result int, searchCount int) {
	mid := len(a) / 2
	val, _ := strconv.Atoi(a[mid])
	switch {
	case len(a) == 0:
		result = -1 // not found
	case val > search:
		result, searchCount = binarySearch(a[:mid], search)
	case val < search:
		result, searchCount = binarySearch(a[mid+1:], search)
		if result >= 0 { // if anything but the -1 "not found" result
			result += mid + 1
		}
	default: // a[mid] == search
		result = mid // found
	}
	searchCount++
	return
}

func fileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func work(client *redis.Client, minioClient *minio.Client) {
	for {
		if queue.GetDeadLen() > 0 {
			if queue.GetDeadLen() > 10 {
				log.Printf("More than 10 messages have failed")
			} else if queue.GetDeadLen() > 100 {
				log.Fatal("More than 100 messages have failed")
			}
			go senders.DoDeadMessages()
		}
		cntx := context.Background()
		res, err := client.BRPop(cntx, 0, "queue").Result()
		if err != nil {
			log.Fatal(err)
		}

		client2 := redis.NewClient(&redis.Options{
			Addr:     "compile-projector:6379",
			Password: "",
			DB:       1,
		})
		_, err_c := client.Ping(cntx).Result()
		if err_c != nil {
			log.Fatal("Unable to connect to queue: ", err_c)
		}

		for i := 0; i < len(res); i++ {
			var job JobCertificate
			err2 := json.Unmarshal([]byte(res[i]), &job)
			if err2 != nil {
				log.Fatal(err2)
			}

			path := filepath.Join(".", fmt.Sprintf("%d", job.BucketID), job.VideoName)
			err3 := os.MkdirAll(path, os.ModePerm)
			if err3 != nil {
				go senders.CompileFailure(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
				continue
			}

			err4 := downloader(minioClient, job.BucketID, job.VideoName, job.FileName)
			if err4 != nil {
				go senders.CompileFailure(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
				continue
			}

			str_bucket_id := fmt.Sprintf("%d", job.BucketID)

			_, not_found := client2.Get(cntx, str_bucket_id+"-"+job.VideoName).Result()
			if not_found == redis.Nil {
				client2.BRPop(cntx, 0, str_bucket_id+"-"+job.VideoName+"-wait").Result()
			} else if not_found != nil {
				log.Panic(not_found)
			}
			client2.Set(cntx, str_bucket_id+"-"+job.VideoName, "1", 0)

			re := regexp.MustCompile("[0-9]+")
			val := re.FindAllString(job.FileName, -1)[0]

			errr := client2.LPush(cntx, str_bucket_id+"-"+job.VideoName+"-done", val).Err()
			if errr != nil {
				go senders.CompileFailure(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
				client2.Del(cntx, str_bucket_id+"-"+job.VideoName)
				client2.LPush(cntx, str_bucket_id+"-"+job.VideoName+"-wait", 1)
				continue
			}
			getList, _ := client2.Sort(cntx, str_bucket_id+"-"+job.VideoName+"-done", &redis.Sort{}).Result()
			int_val, _ := strconv.Atoi(val)
			indx, _ := binarySearch(getList, int_val)

			fileName := fileNameWithoutExtSliceNotation(job.VideoName) + ".gif"

			var cmd *exec.Cmd
			if _, os_err := os.Stat("./" + str_bucket_id + "/" + fileName); errors.Is(os_err, os.ErrNotExist) {
				cmd = exec.Command("convert", "./"+str_bucket_id+"/"+job.VideoName+"/"+job.FileName, "./"+str_bucket_id+"/"+fileName)
			} else if val == getList[0] {
				cmd = exec.Command("convert", "./"+str_bucket_id+"/"+job.VideoName+"/"+job.FileName, "./"+str_bucket_id+"/"+fileName, "./"+str_bucket_id+"/"+fileName)
			} else if val == getList[len(getList)-1] {
				exec.Command("convert", "./"+str_bucket_id+"/"+fileName, "./"+str_bucket_id+"/"+job.VideoName+"/"+job.FileName, "./"+str_bucket_id+"/"+fileName)
			} else {
				exec.Command("convert", "'./"+str_bucket_id+"/"+fileName+"[0"+fmt.Sprintf("%d", indx-1)+"]'", "./"+str_bucket_id+"/"+job.VideoName+"/"+job.FileName, "'./"+str_bucket_id+"/"+fileName+"["+fmt.Sprintf("%d", indx-1)+"--1]'", "./"+str_bucket_id+"/"+fileName)
			}

			err5 := cmd.Run()
			if err5 != nil {
				go senders.CompileFailure(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
				client2.Del(cntx, str_bucket_id+"-"+job.VideoName)
				client2.LPush(cntx, str_bucket_id+"-"+job.VideoName+"-wait", 1)
				continue
			}

			if len(getList) == job.ExpectedFrames {
				err6 := uploader(minioClient, job.BucketID, job.FileName)
				if err6 != nil {
					go senders.CompileFailure(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
					client2.Del(cntx, str_bucket_id+"-"+job.VideoName)
					client2.LPush(cntx, str_bucket_id+"-"+job.VideoName+"-wait", 1)
					continue
				}
			}

			go senders.CompileSuccess(job.BucketID, job.VideoName, job.FileName, job.ExpectedFrames)
			client2.Del(cntx, str_bucket_id+"-"+job.VideoName)
			client2.LPush(cntx, str_bucket_id+"-"+job.VideoName+"-wait", 1)
		}
		client2.Close()
	}
}