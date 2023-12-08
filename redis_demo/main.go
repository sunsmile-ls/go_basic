package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 密码
		DB:       0,  // 数据库
		PoolSize: 20, //连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = rdb.Ping(ctx).Result()
	return err
}

// doCommand go-redis 基本使用
func doCommand() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()

	// 执行命令获取结果
	val, err := rdb.Get(ctx, "key").Result()
	fmt.Println(val, err)

	// 先获取命令对象
	cmder := rdb.Get(ctx, "key")
	fmt.Println(cmder.Val()) // 获取值
	fmt.Println(cmder.Err()) // 获取错误

	// 直接执行命令获取错误
	err = rdb.Set(ctx, "key", 10, time.Hour).Err()
	fmt.Println(err)

	// 直接执行命令获取值
	value := rdb.Get(ctx, "key").Val()
	fmt.Println(value)
}

// doDemo 执行任何命令
func doDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 直接执行命令获取错误
	err := rdb.Do(ctx, "set", "key", 10, "EX", 3600).Err()
	fmt.Println(err)

	// 执行命令获取结果
	val, err := rdb.Do(ctx, "get", "key").Result()
	fmt.Println(val, err)
}

// getValueFormRedis redis.Nil 判断
func getValueFormRedis(key, defaultValue string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		// 如果返回的错误是 key 不存在
		if errors.Is(err, redis.Nil) {
			return defaultValue, nil
		}
		// 出现其他错误
		return "", err
	}
	return val, nil
}

// zsetDemo 操作zset示例
func zsetDemo() {
	// key
	zsetkey := "language_rank"
	// value
	// 注意： v8版本使用[]*redis.Z；此处为v9版本使用 []redis.Z
	language := []redis.Z{
		{Score: 90.0, Member: "Golang"},
		{Score: 98.0, Member: "Java"},
		{Score: 95.0, Member: "Python"},
		{Score: 97.0, Member: "Javascript"},
		{Score: 99.0, Member: "C/C++"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ZADD
	err := rdb.ZAdd(ctx, zsetkey, language...).Err()
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
	fmt.Println("zadd success")

	// 把Golang的分数加10分
	newScore, err := rdb.ZIncrBy(ctx, zsetkey, 10.0, "Golang").Result()
	if err != nil {
		fmt.Printf("ZIncrBy failed, err:%v\n", err)
		return
	}
	fmt.Printf("Golang's score is %f now.\n", newScore)

	// 获取分数最高的3个
	ret := rdb.ZRevRangeWithScores(ctx, zsetkey, 0, 2).Val()
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}

	// 取95 ~ 100分
	op := &redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRangeByScoreWithScores(ctx, zsetkey, op).Result()
	if err != nil {
		fmt.Printf("ZRangeByScoreWithScores failed, err：%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
}

// scanKeysDemo 按前缀查找所有key示例
func scanKeysDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cursor uint64

	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, "prefix:*", 0).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			fmt.Println("key", key)
		}

		if cursor == 0 { // no more keys
			break
		}
	}
}

// scanKeysDemo2 按前缀扫描key
func scanKeysDemo2() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	iter := rdb.Scan(ctx, 0, "prefix:*", 0).Iterator()
	for iter.Next(ctx) {
		fmt.Println("keys", iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
}

// delKeysByMatch 按match格式扫描所有key并删除
func delKeysByMatch(match string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	iter := rdb.Scan(ctx, 0, match, 0).Iterator()
	for iter.Next(ctx) {
		err := rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
}

// pipeline 一次发送多个命令，提高性能，减少往返rtt
func pipeline() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := rdb.Pipeline()
	_, err := p.Set(ctx, "pipeline_counter", 10, time.Hour).Result()
	if err != nil {
		panic(err)
	}
	incr := p.Incr(ctx, "pipeline_counter")

	bc := p.Expire(ctx, "pipeline_counter", time.Hour)
	fmt.Println(bc)

	_, err = p.Exec(ctx)
	if err != nil {
		panic(err)
	}
	// 在执行pipe.Exec之后才能获取到结果
	fmt.Println(incr.Val())
}

// pipelined 一次发送多个命令，提高性能，减少往返rtt
func pipelined() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var incr *redis.IntCmd

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		incr = pipe.Incr(ctx, "pipelined_counter")
		pipe.Expire(ctx, "pipelined_counter", time.Hour)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// 在pipeline执行后获取到结果
	fmt.Println(incr.Val())

}

// tx 事务处理
// Redis 是单线程执行命令的，因此单个命令始终是原子的，但是来自不同客户端的两个给定命令可以依次执行
// 使用事务可以
func tx() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	p := rdb.TxPipeline()
	incr := p.Incr(ctx, "pipelined_counter")
	p.Expire(ctx, "pipelined_counter", time.Hour)
	_, err := p.Exec(ctx)
	fmt.Println(incr.Val(), err)
}

func main() {
	if err := initClient(); err != nil {
		panic(err)
	}
	// doCommand()
	// doDemo()
	// str, err := getValueFormRedis("key1", "sunsmile")
	// fmt.Println(err)
	// fmt.Println(str)

	// zsetDemo()
	// scanKeysDemo()
	// scanKeysDemo2()
	// delKeysByMatch("prefix:*", time.Second*3)
	// pipeline()
	// pipelined()
	tx()
}
