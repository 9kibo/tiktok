package kafka

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"tiktok/biz/dao"
)

var FavoriteMq *TKafka

// ConsumeFavorite 消费逻辑，在协程中异步写入数据库
func ConsumeFavorite(tKafka *TKafka) {
	for {
		msg, err := tKafka.ReadMsg()
		if err != nil {
			logrus.WithField("kafka", msg).Warn("kafka读取失败")
		}
		UserId, _ := strconv.ParseInt(string(msg.Key), 10, 64)
		value := strings.Split(string(msg.Value), " ")
		VideoId, _ := strconv.ParseInt(string(value[0]), 10, 64)
		if len(value) == 1 {
			err = dao.DelFavorite(UserId, VideoId)
		} else {
			createAt, _ := strconv.ParseInt(string(value[1]), 10, 64)
			err = dao.AddFavorite(UserId, VideoId, createAt)
		}
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"UserId":  UserId,
				"VideoId": VideoId,
				"Action":  len(value) ^ 3,
			}).Warn("重复点赞或取消点赞")
		}
	}
}
