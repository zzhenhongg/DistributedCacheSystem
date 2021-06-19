package lru

import "container/list"

//lru 缓存淘汰策略

// LRU Cache
type Cache struct {
	maxBytes  int64                         //最大内存
	nbytes    int64                         //当前已经使用的内存
	ll        *list.List                    //Go的双向链表
	cache     map[string]*list.Element      // key：字符串 value：双向链表的node的指针
	OnEvicted func(key string, value Value) //记录移除时候的回调函数，可以为nil
}

//键值对
type entry struct {
	key   string
	value Value
}

//通过长度len计算需要占用多少的bytes
type Value interface {
	Len() int
}

//实例化Cache的函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//查找功能，1、从字典找到双向链表的指针 2、将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok { //这个鬼东西看起来像是个匿名函数？返回2个值
		c.ll.MoveToFront(ele) //移动到队尾
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//删除功能，淘汰缓存，即淘汰最近最少访问的节点
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() //取得队首的节点（这里back定义为队首）
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                //从字典中删除该节点的映射关系
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) //更新所用内存
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) //调用回调函数
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	//存在则修改并移动到最前面
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	//循环，如果空间不够就一直删除，用for是因为可能删除多个
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 返回添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
