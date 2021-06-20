package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//函数类型Hash 依赖注入
type Hash func(data []byte)uint32

//该map存储所有的哈希键
type Map struct {
	hash Hash
	replicas int //虚拟节点的倍数
	keys []int //哈希环
	hashMap map[int]string //虚拟节点和真实节点的映射表，虚拟节点指向真实节点
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash: fn, //表示自定义的哈希函数
		hashMap: make(map[int]string),
	}
	if m.hashMap == nil {
		m.hash = crc32.ChecksumIEEE //默认的哈希算法
	}
	return m
}

//允许传入一个或者多个真实节点的名称
func (m *Map)Add(keys ...string){
	for _,key := range keys{
		//对于每一个真实节点，创建多个虚拟节点，虚拟节点的名称为strconv.Itoa(i) + key 添加编号
		for i := 0; i < m.replicas;i++{
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys,hash) //将虚拟节点添加到环上
			m.hashMap[hash] = key //添加映射关系
		}
	}
	sort.Ints(m.keys) //环上的哈希值排序
}

func (m *Map)Get(key string)string {
	if len(m.keys) ==0{
		return ""
	}

	hash := int(m.hash([]byte(key))) //计算哈希值
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash //顺时针找到第一批配的虚拟节点的下表idx
	})

	return m.hashMap[m.keys[idx%len(m.keys)]] //因为是环所以取余数
}

