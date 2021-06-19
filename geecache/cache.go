package geecache

import (
	"DistributedCacheSystem/lru"
	"sync"
)

//cache.go 负责并发控制
//实例化 lru，封装 get 和 add 方法，并添加互斥锁 mu
type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//懒汉式初始化 Lazy Initialization
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil) //奇怪，这里cacheBytes并没有设定最大值呀
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
