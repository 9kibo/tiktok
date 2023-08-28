package utils

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"sync"
	"time"
)

//简易全局锁
// 用来保证缓存未命中时，保证相同请求只有一个会打到数据库中

type SafeChan struct {
	sync.Once
	Ch chan struct{}
}

func NewSafeChan() *SafeChan {
	return &SafeChan{
		Ch: make(chan struct{}),
	}
}

func (ch *SafeChan) Close() {
	ch.Do(func() {
		close(ch.Ch)
	})
}

type CacheGuard struct {
	sync.Mutex
	M map[int64]*SafeChan
}

func NewCacheGuard() *CacheGuard {
	return &CacheGuard{
		M: make(map[int64]*SafeChan),
	}
}

// 监听key,如果key未被请求，返回true，否则阻塞直到超时或者redis加载完成
func (CG *CacheGuard) Put(key int64) (bool, error) {
	CG.Lock()
	ch, ok := CG.M[key]
	if !ok {
		ch = NewSafeChan()
		CG.M[key] = ch
		CG.Unlock()
		return true, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	CG.Unlock()
	select {
	case <-ctx.Done():
		delete(CG.M, key)
		return false, errors.New("获取超时")
	case <-ch.Ch:
		return false, nil
	}
}

// 已经获取并加载到redis，通知其余进程
func (CG *CacheGuard) Del(key int64) {
	CG.Lock()
	defer CG.Unlock()
	ch, ok := CG.M[key]
	if !ok {
		return
	}
	ch.Close()
	delete(CG.M, key)
}
