package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"reflect"
	"tiktok/config"
	"time"
	"unsafe"
)

type TKafka struct {
	Topic  string
	Reader *kafka.Reader
	Writer *kafka.Writer
	Ctx    context.Context
}

func GetKafka(Topic string, Group string) *TKafka {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(config.KafkaAddr),
		Topic:                  Topic,
		Balancer:               &kafka.Hash{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: true,
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.KafkaAddr},
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

func (T *TKafka) WriteMsg(key string, value string) {
	T.Writer.WriteMessages(T.Ctx, kafka.Message{
		Key:   String2Bytes(key),
		Value: String2Bytes(value),
	})
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
