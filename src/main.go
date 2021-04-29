package main

import (
	"github.com/Shopify/sarama"
	"github.com/go-ini/ini"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent/src/kafka"
	"logagent/src/tailfile"
	"time"
	"logagent/src/conf"
)

// 日志收集客户端
// tail读文件
// 向kafka中写入数据

// 写入消息的逻辑
func run(tailObj *tail.Tail,config conf.Config) (err error) {
	for{
		// 循环读取数据
		message, ok := <-tailObj.Lines
		if !ok {
			logrus.Warningln("read log err:",err)
			time.Sleep(time.Second)			// 读取失败等待1s
			continue
		}
		logrus.Infoln("message:",message.Text)
		// 利用通道将将同步的代码改为异步的
		// 把读出来的一行日志包装成kafka中的数据类型
		// 构造一个消息
		msg := &sarama.ProducerMessage{}
		msg.Topic = config.Topic
		msg.Value = sarama.StringEncoder(message.Text)
		// 将消息地址写入通道
		kafka.MsgChan <- msg
	}
	return
}

func main() {
	// 0. 读取配置文件
	var config = new(conf.Config)
	err := ini.MapTo(config,"src/conf/conf.ini")
	if err != nil {
		return
	}
	logrus.Infoln("配置文件加载成功",config)

	// 1. 初始化（kakfa连接...）
	client,err := kafka.InitKafka(config)
	if err != nil {
		return
	}

	// 2. 根据配置中的路径使用tail收集日志
	err = tailfile.Init(config.CollectConfig.LogfilePath)
	if err != nil {
		return
	}
	// 3. 通过sarama发到kafka
	// 这里要么两个方法都用goroutine
	// 要么就把run方法放在后面，因为run中有死循环，不这样就执行不到下面的方法
	go kafka.SendMsg(client)
	_ = run(tailfile.TailObj, *config)
}
