package geecache

//缓存值的抽象与封装
//保存一组不可变的bytes 【immutable】
type ByteView struct {
	b []byte //存储真实缓存值 支持任意类型比如字符串、图片
}

func (v ByteView) Len() int {
	return len(v.b)
}

//返回一份复制
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

//克隆一份 缓存只读，防止缓存值被外部程序修改
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b)) //声明一个空间
	copy(c, b)
	return c
}
