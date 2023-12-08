package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8" // 注意导入的是新版本
)

var (
	rdb *redis.Client
)

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 密码
		DB:       0,  // 数据库
		PoolSize: 20, // 连接大小
	})

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	_, err = rdb.Ping(ctx).Result()
	return
}

func doCommand() {
	context.WithTimeout(context.Background(), 50*time.Second)
}
func main() {
	err := initClient()
	if err != nil {
		panic(err)
	}
}
