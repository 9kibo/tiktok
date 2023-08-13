package kafka

//"github.com/segmentio/kafka-go"
//var TestProducer *kafka.Writer

func InitKafka() {
	/*
		TestProducer = &kafka.Writer{
			Addr:                   kafka.TCP(config.KafkaAddr),
			Topic:                  "test",
			Balancer:               &kafka.Hash{},
			WriteTimeout:           1 * time.Second,
			RequiredAcks:           kafka.RequireNone,
			AllowAutoTopicCreation: true,
		}
			Addr            地址
			Topic          指定主题
			Balancer		计算该消息应该被发送到某个分区的方式，默认方法是循环
			MaxAttempts   消息传递的尝试次数
			WriteBackoffMax  消息写入前的最大等待时间 默认1s
			BatchSize	 缓冲区积压消息数量 默认 100 个消息
			BatchBytes	缓冲区大小 默认1048576
			ReadTimeout
			WriteTimeout
			RequiredAcks 写入kafka成功后是否返回ack
			Async  是否开启异步 默认关闭，比较影响性能
			AllowAutoTopicCreation 自动创建 Topic

		调用
	*/
}
