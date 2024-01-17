package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer // 声明一个全局的连接kafka的生产者client
	logDataChan chan *logData
)

// Init 初始化client
func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	// tailf包使⽤
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出⼀个 partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	// 连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	fmt.Println(1)
	// 初始化logDataChan
	logDataChan = make(chan *logData, maxSize)
	// 开启后台的goroutine从通道中取取数据发往kafka
	go sendToKafka()
	return
}

// 给外部暴露的一个函数，该函数只把日志数据发送到一个内部的channel中
func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

// 真正往kafka发送日志的函数
func sendToKafka() {
	for {
		select {
		case ld := <-logDataChan:
			// 构造⼀个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// 发送到kafka
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failed, err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Millisecond * 50)
		}

	}

}

// 优化点：将一条条直接暴露给外部的发数据接口封装成发数据到一个有限的chan，达到限流的目的

// 优化前的代码
// var (
// 	client sarama.SyncProducer
// )

// func Init(addrs []string) (err error) {
// 	config := sarama.NewConfig()

// 	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
// 	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出⼀个partition
// 	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

// 	// 连接kafka producer
// 	client, err = sarama.NewSyncProducer(addrs, config)
// 	if err != nil {
// 		fmt.Println("producer closed, err:", err)
// 		return
// 	}
// 	return
// }

// func SendToKafka(topic, data string) {
// 	// 构造⼀个消息
// 	msg := &sarama.ProducerMessage{}
// 	msg.Topic = topic
// 	msg.Value = sarama.StringEncoder(data)

// 	// 发送消息
// 	pid, offset, err := client.SendMessage(msg)
// 	if err != nil {
// 		fmt.Printf("send msg failed, err: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("pid:%v offset:%v\n", pid, offset)
// }
