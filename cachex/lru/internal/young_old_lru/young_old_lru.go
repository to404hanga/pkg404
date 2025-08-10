package young_old_lru

import (
	"container/list"
	"errors"
	"time"

	"github.com/to404hanga/pkg404/cachex/lru/internal/interfaces"
)

type EvictCallback func(key, value any)

// YoungOldLRU 实现分代LRU缓存算法
type YoungOldLRU struct {
	size          int
	YoungList     *list.List
	OldList       *list.List
	items         map[any]*list.Element
	onEvict       EvictCallback
	youngListSize int
	stayTime      time.Duration
	// 优化：添加时间缓存，减少系统调用
	lastCheckTime time.Time
	checkInterval time.Duration
}

type entry struct {
	key   any
	value any
	addAt time.Time
	flag  bool // true in young, false in old
	// 优化：添加访问计数，用于更智能的晋升策略
	accessCount uint32
}

// NewYoungOldLRU 创建新的YoungOldLRU实例
func NewYoungOldLRU(size int, youngListSize int, stayTime time.Duration, onEvict EvictCallback) (*YoungOldLRU, error) {
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}
	if youngListSize <= 0 {
		return nil, errors.New("must provide a positive youngListSize")
	}
	if youngListSize > size {
		return nil, errors.New("youngListSize must be less than or equal to size")
	}

	return &YoungOldLRU{
		size:          size,
		YoungList:     list.New(),
		OldList:       list.New(),
		items:         make(map[any]*list.Element),
		onEvict:       onEvict,
		youngListSize: youngListSize,
		stayTime:      stayTime,
		// 优化：设置检查间隔，减少时间计算频率
		checkInterval: stayTime / 10, // 每1/10 stayTime检查一次
		lastCheckTime: time.Now(),
	}, nil
}

// Purge 清空缓存
func (c *YoungOldLRU) Purge() {
	for key, value := range c.items {
		if c.onEvict != nil {
			c.onEvict(key, value.Value.(*entry).value)
		}
		delete(c.items, key)
	}
	c.YoungList.Init()
	c.OldList.Init()
}

// Add 添加或更新缓存项
func (c *YoungOldLRU) Add(key, value any) bool {
	if ent, ok := c.items[key]; ok {
		return c.updateExisting(ent, value)
	}

	// 新增元素，直接加入Old队列
	ent := &entry{
		key:         key,
		value:       value,
		addAt:       time.Now(),
		flag:        false,
		accessCount: 1,
	}
	e := c.OldList.PushFront(ent)
	c.items[key] = e

	evicted := c.YoungList.Len()+c.OldList.Len() > c.size
	if evicted {
		c.removeOldest()
	}

	return evicted
}

// updateExisting 更新已存在的缓存项
func (c *YoungOldLRU) updateExisting(ent *list.Element, value any) bool {
	ev := ent.Value.(*entry)
	ev.value = value
	ev.accessCount++

	if ev.flag {
		// 在Young队列中，直接移到前面
		c.YoungList.MoveToFront(ent)
		return false
	}

	// 在Old队列中，检查是否需要晋升
	if c.shouldPromote(ev) {
		c.promoteToYoung(ent)
	} else {
		c.OldList.MoveToFront(ent)
	}
	return false
}

// Get 获取缓存项
func (c *YoungOldLRU) Get(key any) (value any, ok bool) {
	ent, ok := c.items[key]
	if !ok {
		return nil, false
	}

	kv := ent.Value.(*entry)
	kv.accessCount++

	if kv.flag {
		// 快速路径：Young队列命中
		c.YoungList.MoveToFront(ent)
		return kv.value, true
	}

	// Old队列命中，检查是否需要晋升
	if c.shouldPromote(kv) {
		c.promoteToYoung(ent)
	} else {
		c.OldList.MoveToFront(ent)
	}
	return kv.value, true
}

// shouldPromote 判断是否应该晋升到Young队列
func (c *YoungOldLRU) shouldPromote(ev *entry) bool {
	// 优化：减少时间计算频率
	now := time.Now()
	if now.Sub(c.lastCheckTime) > c.checkInterval {
		c.lastCheckTime = now
	}
	
	// 修复：stayTime 是必要条件，访问频率只是加速条件
	timeCondition := ev.addAt.Before(c.lastCheckTime.Add(-c.stayTime))
	if !timeCondition {
		return false // 如果时间不满足，直接返回false
	}
	
	// 时间满足后，可以考虑访问频率作为额外的晋升条件
	// 但在这个版本中，我们保持原有的纯时间逻辑
	return true
}

// promoteToYoung 将元素晋升到Young队列
func (c *YoungOldLRU) promoteToYoung(ent *list.Element) {
	kv := ent.Value.(*entry)
	kv.flag = true
	
	// 修复：先从Old队列移除，再创建新的entry添加到Young队列
	c.OldList.Remove(ent)
	newEnt := c.YoungList.PushFront(kv)
	c.items[kv.key] = newEnt
	
	// 如果Young队列超限，降级最老的元素
	if c.YoungList.Len() > c.youngListSize {
		c.demoteOldestYoung()
	}
}

// demoteOldestYoung 将Young队列中最老的元素降级到Old队列
func (c *YoungOldLRU) demoteOldestYoung() {
	ent := c.YoungList.Back()
	if ent == nil {
		return
	}
	
	kv := ent.Value.(*entry)
	kv.flag = false
	kv.addAt = time.Now()
	
	// 修复：先从Young队列移除，再创建新的entry添加到Old队列
	c.YoungList.Remove(ent)
	newEnt := c.OldList.PushFront(kv)
	c.items[kv.key] = newEnt
}

// Contains 检查key是否存在
func (c *YoungOldLRU) Contains(key any) (ok bool) {
	_, ok = c.items[key]
	return ok
}

// Peek 查看缓存项但不更新位置
func (c *YoungOldLRU) Peek(key any) (value any, ok bool) {
	var ent *list.Element
	if ent, ok = c.items[key]; ok {
		return ent.Value.(*entry).value, true
	}
	return nil, ok
}

// Remove 移除缓存项
func (c *YoungOldLRU) Remove(key any) (present bool) {
	ent, ok := c.items[key]
	if !ok {
		return false
	}

	kv := ent.Value.(*entry)
	if kv.flag {
		c.YoungList.Remove(ent)
	} else {
		c.OldList.Remove(ent)
	}
	delete(c.items, key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
	return true
}

// Resize 调整缓存大小
func (c *YoungOldLRU) Resize(opts ...*interfaces.SizeOptions) (evicted int) {
	var size, youngListSize int
	for _, opt := range opts {
		if opt.Key == "size" {
			size = opt.Value
		}
		if opt.Key == "youngListSize" {
			youngListSize = opt.Value
		}
	}
	if size <= 0 {
		c.Purge()
		return -1
	}
	if youngListSize <= 0 {
		youngListSize = size >> 1
	}

	if youngListSize > size {
		youngListSize = size
	}

	// 如果young队列超过新的限制，将多余元素移到old队列
	diff := c.YoungList.Len() - youngListSize
	for i := 0; i < diff; i++ {
		c.demoteOldestYoung()
	}

	// 如果总长度超过新的size限制，移除最老的元素
	totalLen := c.YoungList.Len() + c.OldList.Len()
	evictCount := 0
	for totalLen > size {
		c.removeOldest()
		evictCount++
		totalLen--
	}

	c.size = size
	c.youngListSize = youngListSize
	return evictCount
}

// RemoveOldest 移除最老的元素
func (c *YoungOldLRU) RemoveOldest() (key, value any, ok bool) {
	ent := c.OldList.Back()
	if ent != nil {
		c.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// GetOldest 获取最老的元素
func (c *YoungOldLRU) GetOldest() (key, value any, ok bool) {
	ent := c.OldList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// Keys 获取所有key
func (c *YoungOldLRU) Keys() []any {
	keys := make([]any, len(c.items))
	i := 0
	for ent := c.OldList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	for ent := c.YoungList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

// Len 获取缓存长度
func (c *YoungOldLRU) Len() int {
	return len(c.items)
}

// removeOldest 移除最老的元素
func (c *YoungOldLRU) removeOldest() {
	ent := c.OldList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

// removeElement 移除指定元素
func (c *YoungOldLRU) removeElement(e *list.Element) {
	kv := e.Value.(*entry)
	if kv.flag {
		c.YoungList.Remove(e)
	} else {
		c.OldList.Remove(e)
	}
	delete(c.items, kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}
