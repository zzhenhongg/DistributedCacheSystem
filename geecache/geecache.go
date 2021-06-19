package geecache

import (
	"fmt"
	"log"
	"sync"
)

//负责与外部交互，控制缓存存储和获取的主流程

type Getter interface {
	Get(key string) ([]byte, error)
}

//实现了Getter接口的方法
type GetterFunc func(key string) ([]byte, error)

//实现接口的一个回调函数 此为接口型函数 由外部决定数据源
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//核心数据结构Group
//一个Group可以理解为是一个缓存的命名空间
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) //存储group的一个字典
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter, //回调函数
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

//返回Group
func GetGroup(name string) *Group {
	mu.RLock() //读锁
	g := groups[name]
	mu.RUnlock()
	return g
}

//从mainCache获取缓存，不存在则加载其他节点的
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

//调用用户回调函数
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key) //调用回调函数
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value) //把数据加入缓存
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
