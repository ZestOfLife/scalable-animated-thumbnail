// server
package main

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB_ADDR = "admin:pass@tcp(127.0.0.1:3306)/status?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	sqlDB, conn_err := gorm.Open(mysql.Open(DB_ADDR), &gorm.Config{})
	if conn_err != nil {
		log.Fatal("Unable to connect to queue", conn_err)
	}

	sqlDB.AutoMigrate(&Store{})
	res := sqlDB.First(&Store{})

	db, close_err := sqlDB.DB()
	if close_err != nil {
		log.Fatal(close_err)
	}
	db.Close()

	cntx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "event-store:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(cntx).Result()
	if err != nil {
		log.Fatal("Unable to connect to event store: ", err)
	}
	log.Println("Connected to event store")

	err = client.XGroupCreate(cntx, "buckets", "observer", "0").Err()
	if err != nil {
		log.Println(err) // Possible that it was created beforehand
	}

	if res.RowsAffected < 1 {
		client.XGroupDestroy(cntx, "buckets", "observer")
		err = client.XGroupCreate(cntx, "buckets", "observer", "0").Err()
		if err != nil {
			log.Println(err) // Possible that it was created beforehand
		}
		streamReader(client, 1)
		log.Println("Updated to current state")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go streamReader(client, 0)
	wg.Wait()
}
