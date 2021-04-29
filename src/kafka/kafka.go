package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"logagent/src/conf"
)

var (
	MsgChan chan *sarama.ProducerMessage
	)

// kafka初始化操作
func InitKafka(config *conf.Config) (client sarama.SyncProducer,err error) {
	logrus.Infoln("开始初始化kafka")

	saconfig := sarama.NewConfig()
	saconfig.Producer.RequiredAcks = sarama.WaitForLocal          // 发送完数据需要leader和follow都确认
	saconfig.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	saconfig.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	logrus.Debugln("config.KafkaConfig.Address:",config.KafkaConfig.Address)
	logrus.Debugln("topic:",config.Topic)
	// 连接kafka
	logrus.Infoln("开始连接kafka")
	client, err = sarama.NewSyncProducer([]string{config.KafkaConfig.Address}, saconfig)
	if err != nil {
		logrus.Error("连接kafka失败, err:", err)
		return
	}
	logrus.Infoln("kafka连接成功")
	// 初始化MsgChan
	MsgChan = make(chan *sarama.ProducerMessage,config.ChanSize)
	return
}

// 从MsgChan中读取msg，发送给kakfa
func SendMsg(client sarama.SyncProducer)  {
	for{
		select {
		case msg := <-MsgChan:
			// 发送消息
			logrus.Infoln("开始发送消息...")
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Errorln("向kafka发送消息失败:", err)
				return
			}
			logrus.Infoln("向kafka发送消息成功")
			logrus.Infoln("pid:", pid,"offset:", offset)
		}
	}
}