package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"reflect"
	"tiktok/biz/config"
	"time"
	"unsafe"
)

type TKafka struct {
	Topic  string
	Reader *kafka.Reader
	Writer *kafka.Writer
	Ctx    context.Context
}

func Init() {
	FavoriteMq = GetKafka("favorite", "favoriteAlter")

	// go 启动消费协程
	go FavConsumer(FavoriteMq)
}
func GetKafka(Topic string, Group string) *TKafka {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(config.C.Kafka.Addr),
		Topic:                  Topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: true,
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.C.Kafka.Addr},
		GroupID: Group,
		Topic:   Topic,
	})
	return &TKafka{
		Topic:  Topic,
		Reader: r,
		Writer: w,
		Ctx:    context.Background(),
	}
}

func (T *TKafka) WriteMsg(key string, value string, back func(string, string)) {
	for i := 0; i < 3; i++ {
		if err := T.Writer.WriteMessages(T.Ctx, kafka.Message{
			Key:   String2Bytes(key),
			Value: String2Bytes(value),
		}); err != nil {
			if err == kafka.LeaderNotAvailable {
				time.Sleep(1 * time.Second)
				continue
			} else {
				back(key, value)
				fmt.Printf("kafka写失败Topic:%s,key:%s,value:%s", T.Topic, key, value)
			}
		} else {
			break
		}
	}
}
func (T *TKafka) ReadMsg() (kafka.Message, error) {
	Msg, err := T.Reader.ReadMessage(T.Ctx)
	if err != nil {
		return kafka.Message{}, err
	}
	return Msg, err
}
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
