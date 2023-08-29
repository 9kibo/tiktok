package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestResult(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6300",
		DB:   0,
	})
	ping := client.Ping(Ctx)
	if ping.Err() != nil {
		panic(ping.Err())
	} else {
		fmt.Println("ping Val", ping.Val())
	}
	//by following result, can see redis result as same sa the redis client by go

	//a is not exists, look the result of not exists kek
	//0 <nil>
	fmt.Println(client.SCard(context.Background(), "a").Result())
	//[] <nil>
	fmt.Println(client.SMembers(context.Background(), "a").Result())
	//0 <nil>
	fmt.Println(client.SRem(context.Background(), "a", 123).Result())

	//add result is add success size
	fmt.Println(client.SAdd(context.Background(), "a", 1, 2, 3).Result())
}
