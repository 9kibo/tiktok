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

	fmt.Println(client.SCard(context.Background(), "a").Result())
	fmt.Println(client.SMembers(context.Background(), "a").Result())
	fmt.Println(client.SRem(context.Background(), "a", 123).Result())
}
