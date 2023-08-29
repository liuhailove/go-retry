package policy

// RtyContextCache 上下文重试cache
type RtyContextCache interface {

	// Get 根据key获取重试上下文
	Get(key interface{}) gretry.RtyContext

	// Put 把Key，ctx加入缓存
	Put(key interface{}, ctx gretry.RtyContext)

	// Remove 从cache中移除
	Remove(key interface{})

	// ContainsKey 判断cache中是否包含key
	ContainsKey(key interface{}) bool
}
