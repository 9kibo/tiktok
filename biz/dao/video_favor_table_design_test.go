package dao

import (
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

// TestVideoFavorCreate 和 TestVideoFavorQuery 证明一下2种结构在性能上没有区别
// 为了方便, 就使用 VideoFavor1 了,因为联合主键支持地比较少
type VideoFavor1 struct {
	Id      int64
	VideoId int64 `gorm:"index"`
	UserId  int64 `gorm:"index"`
}
type VideoFavor2 struct {
	VideoId int64 `gorm:"primaryKey"`
	UserId  int64 `gorm:"primaryKey"`
}

/*
10000
videoFavor1s: 0.073090
videoFavor2s: 0.065169
100000
videoFavor1s: 0.574770
videoFavor2s: 0.417209
*/
func TestVideoFavorCreate(t *testing.T) {
	Db.Migrator().DropTable(&VideoFavor1{}, &VideoFavor2{})
	Db.AutoMigrate(&VideoFavor1{}, &VideoFavor2{})
	size := 100000
	videoFavor1s := make([]*VideoFavor1, 0, size)
	for i := 1; i <= size; i++ {
		videoFavor1s = append(videoFavor1s, &VideoFavor1{
			VideoId: int64(i),
			UserId:  int64(i),
		})
	}
	videoFavor2s := make([]*VideoFavor2, 0, size)
	for i := 1; i <= size; i++ {
		videoFavor2s = append(videoFavor2s, &VideoFavor2{
			VideoId: int64(i),
			UserId:  int64(i),
		})
	}
	bdb := Db.Session(&gorm.Session{CreateBatchSize: 1000})

	start := time.Now()
	r := bdb.Create(&videoFavor1s)
	end := time.Now()
	if r.Error != nil {
		log.Panicln(r.Error)
	}
	if r.RowsAffected != int64(size) {
		log.Panicf("r.RowsAffected[%d] != size[%d]", r.RowsAffected, size)
	}
	log.Printf("videoFavor1s: %f \n", end.Sub(start).Seconds())

	start = time.Now()
	r = bdb.Create(&videoFavor2s)
	end = time.Now()
	if r.Error != nil {
		log.Panicln(bdb.Error)
	}
	if r.RowsAffected != int64(size) {
		log.Panicln(" r.RowsAffected != size ")
	}
	log.Printf("videoFavor2s: %f \n", end.Sub(start).Seconds())
}

/*
data=10000, time=10000
videoFavor1s: 1.093006
videoFavor2s: 1.042756
data=100000, time=100000
videoFavor1s: 9.617810
videoFavor2s: 9.893610
*/
func TestVideoFavorQuery(t *testing.T) {
	size := 10000
	start := time.Now()
	for i := 1; i <= size; i++ {
		r := Db.Find(&VideoFavor1{}, "user_id = ?", i)
		if r.Error != nil {
			t.Log(r.Error)
		}
	}
	end := time.Now()
	log.Printf("videoFavor1s: %f \n", end.Sub(start).Seconds())
	start = time.Now()
	for i := 1; i <= size; i++ {
		r := Db.Find(&VideoFavor2{}, "user_id = ?", i)
		if r.Error != nil {
			t.Log(r.Error)
		}
	}
	end = time.Now()
	log.Printf("videoFavor2s: %f \n", end.Sub(start).Seconds())
}
