package cloudTool

import (
	"fmt"
	"github.com/Shopify/sarama"
)

//InitKaDB 获取kafka生产端链接
func InitKaDB(dsn string) *sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{dsn}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return nil
	}
	defer client.Close()
	return &client
}

//GetKafka 获取消费端
func GetKafka(dsn string) *sarama.Consumer {
	client, err := sarama.NewConsumer([]string{dsn}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return nil
	}
	return &client
}
