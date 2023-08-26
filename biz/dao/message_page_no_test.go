package dao

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"tiktok/biz/model"
	"time"
)

// 查询 用户1 的好友列表 带上最新聊天记录
/* 49个好友
 */
func TestGetFriendList(t *testing.T) {
	type GetFriendListTaskResult struct {
		timeConsuming time.Duration
		result        []*model.Message
		err           error
	}
	type GetFriendListTask struct {
		name          string
		getFriendList func(userId int64, toUserIds []int64) ([]*model.Message, error)
		results       []*GetFriendListTaskResult
		lock          sync.Mutex
	}
	calTimeFunc := func(userId int64, toUserIds []int64, task *GetFriendListTask) {
		start := time.Now()
		result, err := task.getFriendList(userId, toUserIds)
		latency := time.Now().Sub(start)
		task.lock.Lock()
		defer task.lock.Unlock()
		task.results = append(task.results, &GetFriendListTaskResult{
			timeConsuming: latency,
			result:        result,
			err:           err,
		})
	}
	oneTaskTime := 30
	var userId int64 = 1
	pageSize := 50
	toUserIds := make([]int64, 0, pageSize-1)
	for i := 0; i < pageSize-1; i++ {
		toUserIds = append(toUserIds, rand.Int63n(int64(userSize/2))+int64(userSize/10))
	}
	tasks := []*GetFriendListTask{
		{
			name:          "one2AllAndAssociationSubQuery",
			getFriendList: one2AllAndAssociationSubQuery,
		},
		{
			name:          "oneTOne",
			getFriendList: oneTOne,
		},
		{
			name:          "one2AllAndAllMessage",
			getFriendList: one2AllAndAllMessage,
		},
		{
			name:          "one2AllAndAssociationSubQueryConcurrent",
			getFriendList: one2AllAndAssociationSubQueryConcurrent,
		},
		{
			name:          "oneTOneConcurrent",
			getFriendList: oneTOneConcurrent,
		},
		{
			name:          "one2AllAndAllMessageConcurrent",
			getFriendList: one2AllAndAllMessageConcurrent,
		},
	}
	wg := sync.WaitGroup{}
	for _, task := range tasks {
		for i := 0; i < oneTaskTime; i++ {
			wg.Add(1)
			go func(task *GetFriendListTask) {
				defer wg.Done()
				calTimeFunc(userId, toUserIds, task)
			}(task)
		}
	}
	wg.Wait()

	for _, task := range tasks {
		sort.Slice(task.results, func(s2, s1 int) bool {
			return task.results[s2].timeConsuming < task.results[s1].timeConsuming
		})
	}
	sort.Slice(tasks, func(s2, s1 int) bool {
		return tasks[s2].results[0].timeConsuming < tasks[s1].results[0].timeConsuming
	})

	type DataShow struct {
		Desc                      string
		Task                      string
		Err                       error
		timeConsuming             time.Duration
		resultSize                int
		equalWithPageSize         bool
		hasNotMessageForToUserIds []int64
	}
	type Statistics struct {
		Task     string
		Fastest  float64
		Slowest  float64
		Mode     float64
		Median   float64
		Average  float64
		Variance float64
	}
	dataShows := make([]*DataShow, 0, oneTaskTime*len(tasks))
	statisticsList := make([]*Statistics, 0, oneTaskTime*len(tasks))
	for i := 0; i < oneTaskTime; i++ {
		desc := fmt.Sprintf("timeConsuming quick-to-slow-order=%d", i)
		for _, task := range tasks {
			taskResult := task.results[i]
			if taskResult.err != nil {
				dataShows = append(dataShows, &DataShow{
					Desc: desc,
					Task: task.name,
					Err:  taskResult.err,
				})
			} else {
				toUserIdMap := make(map[int64]struct{}, len(taskResult.result))
				for _, result := range taskResult.result {
					if result.ToUserId != userId {
						toUserIdMap[result.ToUserId] = struct{}{}
					} else {
						toUserIdMap[result.FromUserId] = struct{}{}
					}
				}
				hasNotMessageForToUserIds := make([]int64, 0, 0)
				for _, toUserId := range toUserIds {
					if _, ok := toUserIdMap[toUserId]; !ok {
						hasNotMessageForToUserIds = append(hasNotMessageForToUserIds, toUserId)
					}
				}
				dataShows = append(dataShows, &DataShow{
					Desc:                      desc,
					Task:                      task.name,
					timeConsuming:             taskResult.timeConsuming,
					resultSize:                len(taskResult.result),
					equalWithPageSize:         pageSize-1 == len(taskResult.result),
					hasNotMessageForToUserIds: hasNotMessageForToUserIds,
				})
			}
		}
	}

	for _, task := range tasks {
		data := make([]float64, len(task.results))
		for _, result := range task.results {
			data = append(data, float64(result.timeConsuming))
		}
		sort.Float64s(data)
		mode, count := stat.Mode(data, nil)
		if count == 1 {
			mode = 0
		}
		var median float64
		if len(data)%2 == 0 {
			median = (data[len(data)-1/2] + data[len(data)-1/2+1]) / 2
		} else {
			median = data[len(data)-1/2]
		}
		statisticsList = append(statisticsList, &Statistics{
			Task:     task.name,
			Fastest:  floats.Max(data),
			Slowest:  floats.Min(data),
			Mode:     mode,
			Median:   median,
			Average:  stat.Mean(data, nil),
			Variance: stat.Variance(data, nil),
		})
	}

}
func getLatestMessageFromOne2AllAndAll2One(userId int64, userMessageList []*model.Message) []*model.Message {
	sort.Slice(userMessageList, func(s2, s1 int) bool {
		return userMessageList[s2].CreatedAt > userMessageList[s1].CreatedAt
	})

	userLatestMessageMap := make(map[int64]struct{}, len(userMessageList)/2)
	userLatestMessageList := make([]*model.Message, 0, len(userMessageList)/2)
	for _, message := range userMessageList {
		toUserId := message.ToUserId
		if toUserId == userId {
			toUserId = message.FromUserId
		}
		if _, ok := userLatestMessageMap[toUserId]; ok {
			continue
		}
		userLatestMessageMap[toUserId] = struct{}{}
		userLatestMessageList = append(userLatestMessageList, message)
	}
	return userLatestMessageList
}

// one2AllAndAssociationSubQuery
func one2AllAndAssociationSubQuery(userId int64, toUserIds []int64) ([]*model.Message, error) {
	userMessageList := make([]*model.Message, 0, 30)
	tx := Db.Where("from_user_id=? and to_user_id in (?) and id = (?)", userId, toUserIds,
		Db.Select("id").Table("message as m2").Where(
			"from_user_id = ? and to_user_id = message.to_user_id",
			userId).Order("created_at desc").Limit(1),
	).Find(&userMessageList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	toUserMessageList := make([]*model.Message, 0, 30)
	tx = Db.Where("from_user_id in (?) and to_user_id = ? and id = (?)", toUserIds, userId,
		Db.Select("id").Table("message as m2").Where("from_user_id = message.from_user_id  and to_user_id = ?", userId).Order("created_at desc").Limit(1),
	).Find(&toUserMessageList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return getLatestMessageFromOne2AllAndAll2One(userId, append(userMessageList, toUserMessageList...)), nil
}
func oneTOne(userId int64, toUserIds []int64) ([]*model.Message, error) {
	messageList := make([]*model.Message, 0, 30)
	for _, toUserId := range toUserIds {
		message := model.Message{}
		tx := Db.Where("(from_user_id = ? and to_user_id =?) or (to_user_id = ? and from_user_id =?)",
			userId, toUserId, toUserId, userId,
		).Order("created_at desc").Take(&message)
		if tx.Error != nil {
			return nil, tx.Error
		}
		messageList = append(messageList, &message)
	}
	return messageList, nil
}
func one2AllAndAllMessage(userId int64, toUserIds []int64) ([]*model.Message, error) {
	filterMap := make(map[int64]struct{}, 30)
	toSelectMessageList := make([]*model.Message, 0, 30)
	userMessageList := make([]*model.Message, 0, 30)
	tx := Db.Table("(?) as latest",
		Db.Model(&model.Message{}).Where("from_user_id = ? and to_user_id in (?)", userId, toUserIds).Order("to_user_id, created_at desc"),
	).Find(&userMessageList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	for _, message := range userMessageList {
		if _, ok := filterMap[message.ToUserId]; ok {
			continue
		}
		filterMap[message.ToUserId] = struct{}{}
		toSelectMessageList = append(toSelectMessageList, message)
	}

	toUserMessageList := make([]*model.Message, 0, 30)
	tx = Db.Table("(?) as latest",
		Db.Model(&model.Message{}).Where("from_user_id in (?) and to_user_id = ?", toUserIds, userId).Order("from_user_id, created_at desc"),
	).Find(&toUserMessageList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	for _, message := range toUserMessageList {
		if _, ok := filterMap[message.FromUserId]; ok {
			continue
		}
		filterMap[message.FromUserId] = struct{}{}
		toSelectMessageList = append(toSelectMessageList, message)
	}

	return getLatestMessageFromOne2AllAndAll2One(userId, toSelectMessageList), nil
}

// one2AllAndAssociationSubQuery
func one2AllAndAssociationSubQueryConcurrent(userId int64, toUserIds []int64) ([]*model.Message, error) {
	var err error = nil
	wg := sync.WaitGroup{}
	userMessageList := make([]*model.Message, 0, 30)
	toUserMessageList := make([]*model.Message, 0, 30)

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		tx := Db.WithContext(ctx).Where("from_user_id=? and to_user_id in (?) and id = (?)", userId, toUserIds,
			Db.Select("id").Table("message as m2").Where(
				"from_user_id = ? and to_user_id = message.to_user_id",
				userId).Order("created_at desc").Limit(1),
		).Find(&userMessageList)
		if tx.Error != nil {
			cancelFunc()
			err = fmt.Errorf("when get user and friends chat record, err=%s", tx.Error)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		tx := Db.WithContext(ctx).Where("from_user_id in (?) and to_user_id = ? and id = (?)", toUserIds, userId,
			Db.Select("id").Table("message as m2").Where("from_user_id = message.from_user_id  and to_user_id = ?",
				userId).Order("created_at desc").Limit(1),
		).Find(&toUserMessageList)
		if tx.Error != nil {
			cancelFunc()
			err = fmt.Errorf("when get friends and user chat record, err=%s", tx.Error)
			return
		}
	}()
	wg.Wait()
	return getLatestMessageFromOne2AllAndAll2One(userId, append(userMessageList, toUserMessageList...)), err
}

func oneTOneConcurrent(userId int64, toUserIds []int64) ([]*model.Message, error) {
	var err error = nil
	messageList := make([]*model.Message, 0, 30)
	wg := sync.WaitGroup{}
	ctx, cancelFunc := context.WithCancel(context.Background())
	for _, toUserId := range toUserIds {
		wg.Add(1)
		go func(toUserId int64) {
			defer wg.Done()

			message := model.Message{}
			tx := Db.WithContext(ctx).Where("(from_user_id = ? and to_user_id =?) or (to_user_id = ? and from_user_id =?)",
				userId, toUserId, toUserId, userId,
			).Order("created_at desc").Take(&message)
			if tx.Error != nil {
				cancelFunc()
				err = fmt.Errorf("when get user[%d] and toUser[%d] chat record, err=%s", userId, toUserId, tx.Error)
				return
			}
			messageList = append(messageList, &message)
		}(toUserId)
	}
	wg.Wait()
	return messageList, err
}
func one2AllAndAllMessageConcurrent(userId int64, toUserIds []int64) ([]*model.Message, error) {
	var err error = nil
	filterMap := make(map[int64]struct{}, 30)
	toSelectMessageList := make([]*model.Message, 0, 30)
	userMessageList := make([]*model.Message, 0, 30)
	toUserMessageList := make([]*model.Message, 0, 30)
	wg := sync.WaitGroup{}
	ctx, cancelFunc := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		tx := Db.WithContext(ctx).Table("(?) as latest",
			Db.Model(&model.Message{}).Where("from_user_id = ? and to_user_id in (?)", userId, toUserIds).Order("to_user_id, created_at desc"),
		).Find(&userMessageList)
		if tx.Error != nil {
			cancelFunc()
			err = fmt.Errorf("when get user and friends chat record, err=%s", tx.Error)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tx := Db.WithContext(ctx).Table("(?) as latest",
			Db.Model(&model.Message{}).Where("from_user_id in (?) and to_user_id = ?", toUserIds, userId).Order("from_user_id, created_at desc"),
		).Find(&toUserMessageList)
		if tx.Error != nil {
			cancelFunc()
			err = fmt.Errorf("when get friends and user  chat record, err=%s", tx.Error)
		}
	}()
	wg.Wait()
	for _, message := range userMessageList {
		if _, ok := filterMap[message.ToUserId]; ok {
			continue
		}
		filterMap[message.ToUserId] = struct{}{}
		toSelectMessageList = append(toSelectMessageList, message)
	}
	for _, message := range toUserMessageList {
		if _, ok := filterMap[message.FromUserId]; ok {
			continue
		}
		filterMap[message.FromUserId] = struct{}{}
		toSelectMessageList = append(toSelectMessageList, message)
	}
	return getLatestMessageFromOne2AllAndAll2One(userId, toSelectMessageList), err
}

var (
	userSize = 10000
)

// 添加 userSize-1 个用户 与 用户1 的对话, 每个用户给用户1发生messageSize个消息, 用户1也是
func TestMessageDataAdd(t *testing.T) {
	var (
		userId         int64 = 1
		messageSize          = 30
		createdAtStart int64 = 1692000000
		createdAtRange       = int64(messageSize) * 1000
	)

	messages := make([]*model.Message, 0, (userSize-1)*messageSize*2)

	for toUserId := 2; toUserId <= userSize; toUserId++ {
		for i := 0; i < messageSize; i++ {
			messages = append(messages, &model.Message{
				FromUserId: userId,
				ToUserId:   int64(toUserId),
				Content:    fmt.Sprintf("from user[%d] to user[%d]", userId, toUserId),
				CreatedAt:  createdAtStart + rand.Int63n(createdAtRange),
			})
		}
		for i := 0; i < messageSize; i++ {
			messages = append(messages, &model.Message{
				FromUserId: int64(toUserId),
				ToUserId:   userId,
				Content:    fmt.Sprintf("from user[%d] to user[%d]", toUserId, userId),
				CreatedAt:  createdAtStart + rand.Int63n(createdAtRange),
			})
		}
	}
	tx := Db.CreateInBatches(&messages, 1000)
	assert.NoError(t, tx.Error)
	assert.Equal(t, int64(len(messages)), tx.RowsAffected)
}
