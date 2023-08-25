package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"tiktok/biz/dao"
	"tiktok/biz/model"
)

var CommonMq *TKafka

func ConsumeComm(tKafka *TKafka) {
	for {
		msg, err := tKafka.ReadMsg()
		if err != nil {
			logrus.WithField("kafka", msg).Warn("kafka读取失败")
		}
		var CommInfo model.CommToJson
		err = json.Unmarshal(msg.Value, &CommInfo)
		if err != nil {
			logrus.WithField("kafka", msg).WithField("Info", fmt.Sprintf("%v", CommInfo)).Warn("comment解析失败")
		}
		if Bytes2String(msg.Key) == "Add" {
			err = dao.AddComm(&model.Comment{
				Id:        CommInfo.CommId,
				CreatedAt: CommInfo.CreateAt,
				UserId:    CommInfo.UserId,
				VideoId:   CommInfo.VideoId,
				Content:   CommInfo.Content,
			})
			if err != nil {
				logrus.WithField("kafka", msg).WithField("Info", fmt.Sprintf("%v", CommInfo)).Warn("添加评论失败")
			}
		} else {
			err = dao.DelComm(CommInfo.CommId, CommInfo.UserId)
			if err != nil {
				logrus.WithField("kafka", msg).WithField("Info", fmt.Sprintf("%v", CommInfo)).Warn("删除评论失败")
			}
		}
	}
}
