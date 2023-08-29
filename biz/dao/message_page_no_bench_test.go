package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

/*
*
测试:

	1.基准测试: 使用下面命令把基准测试结果输出到一个文件中
		go test -v -count 30 -benchtime 30x -benchmem -run=none -bench .> ../../statistics/GetFriendList/bench/all.txt
	2.调用 TestCollectBenchmarkData 生成json数据, 可以选择其他工具或者go的excelize把json转为excel
	3.把excel上传https://dycharts.com/ 生成折线图进行分析

	4.基准测试仅仅是执行30次函数的执行平均耗时, 平均分配内存分析(每个方法有30组数据)
	5.如果需要每个函数的执行时间的的分析, 执行 TestGetFriendList, statistics/GetFriendList/data生成 samplingFrequency 个
		json文件, 每个方法有 taskFuncRunTimes 组数据(有单位是mircos, ms的)
		转excel是一样的方法
		statistics/GetFriendList/linechart是选前5个进行生成折线图

经分析折线图和benchmark
 1. 在耗时上, OneTOneConcurrent是相对快的(4ms-5ms)
    其次是One2AllAndAllMessageConcurrent(7ms-10ms),One2AllAndAllMessage(10ms-12ms).
    OneTOne(15ms-28ms)与上面几个相差有一点点大
    但是One2AllAndAssociationSubQuery和One2AllAndAssociationSubQueryConcurrent是真的特别慢,在60ms-100ms左右
 2. 在内存分配上,
    One2AllAndAssociationSubQuery和One2AllAndAssociationSubQueryConcurrent是最小的(100000B左右)
    OneTOne, OneTOneConcurrent居中(500000B-600000B左右)
    One2AllAndAllMessageConcurrent,One2AllAndAllMessage是最大的(7000000B左右)

原因:
1. One2AllAndAssociationSubQuery涉及关联子查询, 但查的数据最少, 仅仅是用户对好友和好友对用户的最新消息
2. One2AllAndAllMessage仅仅是返回不相关子查询, 但是数据特别多, 双方消息全返回
3. OneTOne是轮询每个好友进行查询直接拿到最新消息, 因此在速度上比One2AllAndAllMessage慢一些, 但是并发的话就非常快

结论: 采用OneTOneConcurrent
*/
var (
	userId, toUserIds, pageSize = GetArgs()
)

func BenchmarkOne2AllAndAssociationSubQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messages, err := one2AllAndAssociationSubQuery(userId, toUserIds)
		if err != nil {
			b.Error(err)
		}
		b.StopTimer()
		assertTrue(userId, toUserIds, messages)
		b.StartTimer()
	}
}
func BenchmarkOneTOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messages, err := oneTOne(userId, toUserIds)
		if err != nil {
			b.Error(err)
		}
		b.StopTimer()
		assertTrue(userId, toUserIds, messages)
		b.StartTimer()
	}
}
func BenchmarkOne2AllAndAllMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messages, err := one2AllAndAllMessage(userId, toUserIds)
		if err != nil {
			b.Error(err)
		}
		b.StopTimer()
		assertTrue(userId, toUserIds, messages)
		b.StartTimer()
	}
}
func BenchmarkOne2AllAndAssociationSubQueryConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messages, err := one2AllAndAssociationSubQueryConcurrent(userId, toUserIds)
		if err != nil {
			b.Error(err)
		}
		b.StopTimer()
		assertTrue(userId, toUserIds, messages)
		b.StartTimer()
	}
}
func BenchmarkOneTOneConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messages, err := oneTOneConcurrent(userId, toUserIds)
		if err != nil {
			b.Error(err)
		}
		b.StopTimer()
		assertTrue(userId, toUserIds, messages)
		b.StartTimer()
	}
}
func BenchmarkOne2AllAndAllMessageConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		messages, err := one2AllAndAllMessageConcurrent(userId, toUserIds)
		if err != nil {
			b.Error(err)
		}
		b.StopTimer()
		assertTrue(userId, toUserIds, messages)
		b.StartTimer()
	}
}

type BenchmarkTask struct {
	Name   string
	Result []*BenchmarkResult
	order  int
}
type BenchmarkResult struct {
	Order     int
	N         int64
	Bytes     int64
	MemAllocs int64
	Count     int64
}

// 从一个文件(文件名无所谓)中解析benchmark
/**
规范:
1.一个文件可以有多个基准测试, 每个基准测试会生成一个json文件
*/
/**
取值
1. Task
Name=BenchmarkOneTOneConcurrent-8除去Benchmark
order是指一个文件可以有多个基准测试的情况下种基准测试的顺序
2. Result
Order是指基准测试指定了count参数, 这种基准测试执行的顺序
Count,N,Bytes,MemAllocs 分别取与从测试名后开始的第2,3,4,5个参数值
	Bytes,MemAllocs需要指定参数-benchmem
*/
// 删除多余的日志, 只保留结果日志
func getBenchmarkTasks(path string) ([]*BenchmarkTask, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")
	btMap := make(map[string]*BenchmarkTask)
	var curTask *BenchmarkTask
	taskOrder := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "Benchmark") || !strings.Contains(line, "/op") {
			continue
		}
		result := strings.Split(line, "\t")
		if len(result) < 5 {
			continue
		}

		ri := 1
		br := BenchmarkResult{}
		for _, s := range result {
			s = strings.TrimSpace(s)
			if s != "" {
				if ri == 1 {
					if !strings.HasPrefix(s, "Benchmark") {
						return nil, errors.New("must a benchmark test func name at the line 1th")
					}
					btName := strings.Replace(s, "Benchmark", "", 1)
					if curTask = btMap[btName]; curTask == nil {
						curTask = &BenchmarkTask{
							Name:  btName,
							order: taskOrder,
						}
						taskOrder++
						btMap[btName] = curTask
					}
				} else if ri == 2 {
					br.Count, err = strconv.ParseInt(s, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("N must a int at the line 2th value, %e", err)
					}
				} else {
					metric := strings.Split(s, " ")
					intV, err := strconv.ParseInt(metric[0], 10, 64)
					if err != nil {
						return nil, fmt.Errorf("must a int at the line 3/4/5th value, %e", err)
					}
					if metric[1] == "ns/op" {
						br.N = intV
					} else if metric[1] == "B/op" {
						br.Bytes = intV
					} else if metric[1] == "allocs/op" {
						br.MemAllocs = intV
					}
				}
				ri++
			}
		}
		br.Order = len(curTask.Result) + 1
		curTask.Result = append(curTask.Result, &br)
	}
	bts := make([]*BenchmarkTask, 0, len(btMap))
	for _, task := range btMap {
		bts = append(bts, task)
	}
	sort.Slice(bts, func(s2, s1 int) bool {
		return bts[s1].order > bts[s2].order
	})
	return bts, nil
}
func TestCollectBenchmarkData(t *testing.T) {
	benchmarkOutput := "A:\\code\\backend\\go\\tiktok\\statistics\\GetFriendList\\bench\\all.txt"
	bts, err := getBenchmarkTasks(benchmarkOutput)
	if err != nil {
		t.Error(err)
	}
	dir := filepath.Dir(benchmarkOutput)
	writeJSON := func(data any, path string) {
		marshal, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path, marshal, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}
	//for _, bt := range bts {
	//	writeJSON(&bt, filepath.Join(dir, bt.Name+".json"))
	//}
	rlen := 0
	for _, bt := range bts {
		if rlen == 0 {
			rlen = len(bt.Result)
		} else if rlen != len(bt.Result) {
			t.Fatal("result len not equal")
		}
	}
	nMapList := make([]map[string]any, 0, rlen)
	bytesMapList := make([]map[string]any, 0, rlen)
	mAMapList := make([]map[string]any, 0, rlen)

	for _, bt := range bts {
		sort.Slice(bt.Result, func(s2, s1 int) bool {
			return bt.Result[s1].N > bt.Result[s2].N
		})
	}
	for i := 0; i < rlen; i++ {
		nMap := make(map[string]any, 7)
		nMap["X"] = i + 1
		for _, bt := range bts {
			nMap[bt.Name] = bt.Result[i].N
		}
		nMapList = append(nMapList, nMap)
	}

	for _, bt := range bts {
		sort.Slice(bt.Result, func(s2, s1 int) bool {
			return bt.Result[s1].Bytes > bt.Result[s2].Bytes
		})
	}
	for i := 0; i < rlen; i++ {
		bytesMap := make(map[string]any, 7)
		bytesMap["X"] = i + 1
		for _, bt := range bts {
			bytesMap[bt.Name] = bt.Result[i].Bytes
		}
		bytesMapList = append(bytesMapList, bytesMap)
	}

	for _, bt := range bts {
		sort.Slice(bt.Result, func(s2, s1 int) bool {
			return bt.Result[s1].MemAllocs > bt.Result[s2].MemAllocs
		})
	}
	for i := 0; i < rlen; i++ {
		mAMap := make(map[string]any, 7)
		mAMap["X"] = i + 1
		for _, bt := range bts {
			mAMap[bt.Name] = bt.Result[i].MemAllocs
		}
		mAMapList = append(mAMapList, mAMap)
	}

	writeJSON(&nMapList, filepath.Join(dir, "N.json"))
	writeJSON(&bytesMapList, filepath.Join(dir, "bytes.json"))
	writeJSON(&mAMapList, filepath.Join(dir, "MemAllocs.json"))
}
