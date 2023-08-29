package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"
	"tiktok/biz/model"
	"time"
)

// 查询 用户1 的好友列表 带上最新聊天记录
/* 49个好友
 */
type Statistics struct {
	Fastest  float64
	Slowest  float64
	Mode     float64
	Median   float64
	Average  float64
	Variance float64
}
type GetFriendListTask struct {
	Name                      string
	result                    []*model.Message
	timeConsuming             time.Duration
	TimeConsuming             string
	Err                       error
	ResultSize                int
	RightResultSize           bool
	HasNotMessageForToUserIds []int64
}
type TaskGroup struct {
	lock          sync.Mutex
	getFriendList func(userId int64, toUserIds []int64) ([]*model.Message, error)
	name          string
	taskList      []*GetFriendListTask
	Statistics    *Statistics
}

// taskList already sorted
func getStatistics(taskList []*GetFriendListTask) *Statistics {
	data := make([]float64, 0, len(taskList))
	for _, result := range taskList {
		data = append(data, float64(result.timeConsuming.Truncate(time.Microsecond)))
	}
	mode, count := stat.Mode(data, nil)
	if count == 1 {
		mode = 0
	}
	var median float64
	dataLen := len(data)
	if dataLen%2 == 0 {
		median = (data[dataLen/2] + data[dataLen/2+1]) / 2
	} else {
		median = data[dataLen/2+1]
	}
	return &Statistics{
		Fastest:  floats.Min(data),
		Slowest:  floats.Max(data),
		Mode:     mode,
		Median:   median,
		Average:  stat.Mean(data, nil),
		Variance: stat.Variance(data, nil),
	}
}

var taskGroups = []*TaskGroup{
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

// TestGetFriendList first to run TestMessageDataAdd add test data to db
// 折线图分析 https://dycharts.com/ 需要把json转excel --使用excelize
/**
#one2AllAndAssociationSubQuery 查询2次, 有关联子查询
SELECT *
FROM `message`
WHERE from_user_id = 1
  and to_user_id in
      (5410, 4551)
  and id = (SELECT id
            FROM message as m2
            WHERE from_user_id = 1
              and to_user_id = message.to_user_id
            ORDER BY created_at desc
            LIMIT 1);
SELECT *
FROM `message`
WHERE from_user_id in
      (5410, 4551)
  and to_user_id = 1
  and id = (SELECT id
            FROM message as m2
            WHERE from_user_id = message.from_user_id
              and to_user_id = 1
            ORDER BY created_at desc
            LIMIT 1);

#one2AllAndAllMessage 查询2次, 无关联子查询, 但是会把一对多用户的所有消息查出来
SELECT *
FROM (SELECT *
      FROM `message`
      WHERE from_user_id = 1
        and to_user_id in
            (5410, 4551)
      ORDER BY to_user_id, created_at desc) as latest;
SELECT *
FROM (SELECT *
      FROM `message`
      WHERE from_user_id in
            (5410, 4551)
        and to_user_id = 1
      ORDER BY from_user_id, created_at desc) as latest;

#oneTOne 一对一地查询
SELECT *
FROM `message`
WHERE (from_user_id = 1 and to_user_id = 5410)
   or (to_user_id = 5410 and from_user_id = 1)
ORDER BY created_at desc
LIMIT 1
*/
func TestGetFriendList(t *testing.T) {
	samplingFrequency := 20
	dataDir := "A:\\code\\backend\\go\\tiktok\\statistics"
	//every taskGroup func run times, range is any
	taskFuncRunTimes := 30
	//assume the current user is id = 1

	theTask := func(samplingIndex int) {
		//do taskGroup for 6 func
		wg := sync.WaitGroup{}
		for _, taskGroup := range taskGroups {
			for i := 0; i < taskFuncRunTimes; i++ {
				wg.Add(1)
				go func(task *TaskGroup) {
					defer wg.Done()
					start := time.Now()
					result, err := task.getFriendList(userId, toUserIds)
					latency := time.Now().Sub(start)
					task.lock.Lock()
					defer task.lock.Unlock()
					task.taskList = append(task.taskList, &GetFriendListTask{
						Name:          task.name,
						result:        result,
						timeConsuming: latency,
						Err:           err,
					})
				}(taskGroup)
			}
			wg.Wait()
		}

		//sort taskGroup.taskList by timeConsuming aes
		for _, taskGroup := range taskGroups {
			sort.Slice(taskGroup.taskList, func(s2, s1 int) bool {
				return taskGroup.taskList[s2].timeConsuming < taskGroup.taskList[s1].timeConsuming
			})
		}
		//sort taskGroup by result.fastest aes
		sort.Slice(taskGroups, func(s2, s1 int) bool {
			return taskGroups[s2].taskList[0].timeConsuming < taskGroups[s1].taskList[0].timeConsuming
		})

		for _, taskGroup := range taskGroups {
			//check the taskList, every taskGroup'result has right result size and user the last message for send to friend
			//	must have a message for not is not friend
			for _, taskResult := range taskGroup.taskList {
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
				taskResult.ResultSize = len(taskResult.result)
				taskResult.RightResultSize = pageSize-1 == len(taskResult.result)
				taskResult.HasNotMessageForToUserIds = hasNotMessageForToUserIds
				taskResult.TimeConsuming = taskResult.timeConsuming.String()
			}
			//statistics:
			//taskGroup.Statistics = getStatistics(taskGroup.taskList)
		}

		//write data
		lineChartDataMicroSs := make([]map[string]any, 0, len(taskGroups)*taskFuncRunTimes)
		lineChartDataMSs := make([]map[string]any, 0, len(taskGroups)*taskFuncRunTimes)
		for i := 0; i < taskFuncRunTimes; i++ {
			lineChartDataMicroS := make(map[string]any, 7)
			lineChartDataMicroS["X"] = i + 1
			lineChartDataMS := make(map[string]any, 7)
			lineChartDataMS["X"] = i + 1
			for _, taskGroup := range taskGroups {
				lineChartDataMicroS[taskGroup.name] = int64(taskGroup.taskList[i].timeConsuming.Microseconds())
				lineChartDataMS[taskGroup.name] = int64(taskGroup.taskList[i].timeConsuming.Milliseconds())
			}
			lineChartDataMicroSs = append(lineChartDataMicroSs, lineChartDataMicroS)
			lineChartDataMSs = append(lineChartDataMSs, lineChartDataMS)
		}
		writeJson := func(path string, data any) {
			marshal, err := json.Marshal(data)
			if err != nil {
				t.Fatal(err)
			}
			err = os.WriteFile(path, marshal, os.ModePerm)
			if err != nil {
				t.Fatal(err)
			}
		}
		writeJson(filepath.Join(dataDir, fmt.Sprintf("test%d-micros.json", samplingIndex)), &lineChartDataMicroSs)
		writeJson(filepath.Join(dataDir, fmt.Sprintf("test%d-ms.json", samplingIndex)), &lineChartDataMSs)
	}

	for i := 0; i < samplingFrequency; i++ {
		theTask(i + 1)
	}

	//t.Log("logging, time unit for float is time.Microsecond")
	////print statistics
	//for _, taskGroup := range taskGroups {
	//	marshal, err := json.MarshalIndent(taskGroup.Statistics, " ", "  ")
	//	if err != nil {
	//		t.Errorf("taskGroup=%s\n\tjson err=%s", taskGroup.name, err)
	//	}
	//	t.Logf("taskGroup=%s, Statistics=\n%s\n%s", taskGroup.name, string(marshal), strings.Repeat("=", 20))
	//}
	//
	////print taskList one to one group
	//for i := 0; i < taskFuncRunTimes; i++ {
	//	taskOneToOneGroup := make([]*GetFriendListTask, 0, len(taskGroups))
	//	for _, taskGroup := range taskGroups {
	//		taskOneToOneGroup = append(taskOneToOneGroup, taskGroup.taskList[i])
	//	}
	//	marshal, err := json.MarshalIndent(taskOneToOneGroup, " ", "  ")
	//	if err != nil {
	//		t.Errorf("the taskOneToOneGroup=%d\n\tjson err=%s", i+1, err)
	//	}
	//	t.Logf("the taskOneToOneGroup=%d=%s\n%s", i+1, string(marshal), strings.Repeat("=", 20))
	//}

}

func GetArgs() (int64, []int64, int) {
	var userId int64 = 1
	//assume user friend id range is [2,pageSize-1]
	pageSize := 50
	toUserIds := make([]int64, 0, pageSize-1)
	for i := 0; i < pageSize-1; i++ {
		toUserIds = append(toUserIds, rand.Int63n(int64(userSize/2))+int64(userSize/10))
	}
	return userId, toUserIds, pageSize
}
func assertTrue(userId int64, toUserIds []int64, messageList []*model.Message) {
	toUserIdMap := make(map[int64]struct{}, len(messageList))
	for _, result := range messageList {
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
	if len(hasNotMessageForToUserIds) != 0 {
		panic(fmt.Sprintf("fail, for hasNotMessageForToUserIds=%v", hasNotMessageForToUserIds))
	}
}
func TestOne2AllAndAssociationSubQuery(t *testing.T) {
	messages, err := one2AllAndAssociationSubQuery(userId, toUserIds)
	if err != nil {
		t.Fatal(err)
	}
	assertTrue(userId, toUserIds, messages)
}
func TestOneTOne(t *testing.T) {
	messages, err := oneTOne(userId, toUserIds)
	if err != nil {
		t.Fatal(err)
	}
	assertTrue(userId, toUserIds, messages)
}
func TestOne2AllAndAllMessage(t *testing.T) {
	messages, err := one2AllAndAllMessage(userId, toUserIds)
	if err != nil {
		t.Fatal(err)
	}
	assertTrue(userId, toUserIds, messages)
}
func TestOne2AllAndAssociationSubQueryConcurrent(t *testing.T) {
	messages, err := one2AllAndAssociationSubQueryConcurrent(userId, toUserIds)
	if err != nil {
		t.Fatal(err)
	}
	assertTrue(userId, toUserIds, messages)
}
func TestOneTOneConcurrent(t *testing.T) {
	messages, err := oneTOneConcurrent(userId, toUserIds)
	if err != nil {
		t.Fatal(err)
	}
	assertTrue(userId, toUserIds, messages)
}
func TestOne2AllAndAllMessageConcurrent(t *testing.T) {
	messages, err := one2AllAndAllMessageConcurrent(userId, toUserIds)
	if err != nil {
		t.Fatal(err)
	}
	assertTrue(userId, toUserIds, messages)
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
	wg.Add(2)
	userMessageList := make([]*model.Message, 0, 30)
	toUserMessageList := make([]*model.Message, 0, 30)

	ctx, cancelFunc := context.WithCancel(context.Background())
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
	messageList := make([]*model.Message, len(toUserIds))
	wg := sync.WaitGroup{}
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg.Add(len(toUserIds))
	for i, toUserId := range toUserIds {
		go func(i int, toUserId int64) {
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
			messageList[i] = &message
		}(i, toUserId)
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
