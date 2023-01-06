package db

import (
	"github.com/Shopify/sarama"
	"github.com/laydong/toolpkg"
	"github.com/laydong/toolpkg/logx"
)

/**
InitKafkaProducer 获取kafka生产端
dsn string localhost:9093
username 账号 可传空
password 密码
*/
func InitKafkaProducer(dsn, username, password string) (db *sarama.SyncProducer, err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	if username != "" && password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
	}
	// 连接kafka
	client, er := sarama.NewSyncProducer([]string{dsn}, config)
	if er != nil {
		logx.ErrorF(toolpkg.GetNewGinContext(), "kafka生产者链接错误", er.Error())
		return
	}
	defer client.Close()
	return &client, er
}

/**InitKafkaConsumer 获取消费端
dsn string localhost:9093
username 账号 可传空
password 密码
*/
func InitKafkaConsumer(dsn, username, password string) (db *sarama.Consumer, err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	if username != "" && password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
	}
	// 连接kafka
	client, err := sarama.NewConsumer([]string{dsn}, config)
	if err != nil {
		logx.ErrorF(toolpkg.GetNewGinContext(), "kafka生产者链接错误", err.Error())
		return
	}
	return &client, err
}
