package tailfile

import (
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var (
	TailObj *tail.Tail
)

func Init(filename string) (err error) {
	config := tail.Config{
		ReOpen: true,									// 重新打开
		Follow: true,									// 是否跟随
		Location: &tail.SeekInfo{Offset: 0,Whence: 2},	// 从文件的哪个地方开始读
		MustExist: false,								// 文件不存在不报错
		Poll: true,
	}

	// 打开文件开始读取
	logrus.Infoln("初始化tail操作")
	TailObj, err = tail.TailFile(filename, config)
	if err != nil {
		logrus.Errorln("tail操作打开文件失败:",err)
		return
	}
	logrus.Infoln("tail操作初始化成功")
	return err
}
