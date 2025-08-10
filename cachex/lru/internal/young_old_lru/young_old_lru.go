package young_old_lru

import (
	"container/list"
	"errors"
	"time"

	"github.com/to404hanga/pkg404/cachex/lru/internal/interfaces"
)

type EvictCallback func(key, value any)

type YoungOldLRU struct {
	size          int
	YoungList     *list.List
	OldList       *list.List
	items         map[any]*list.Element
	onEvict       EvictCallback
	youngListSize int
	stayTime      time.Duration
}

type entry struct {
	key   any
	value any
	addAt time.Time
	flag  bool // true in young, false in old
}

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
	}, nil
}

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

func (c *YoungOldLRU) Add(key, value any) bool {
	if ent, ok := c.items[key]; ok {
		ev := ent.Value.(*entry)
		if ev.flag {
			// young 队列
			c.YoungList.MoveToFront(ent)
			return false
		} else {
			// old 队列
			if ev.addAt.Before(time.Now().Add(-c.stayTime)) {
				// 如果在 old 队列中的时间超过了 stayTime，允许晋升到 young 队列
				c.promote(key)
				return false
			}
			// 否则只移动到 old 队列头部
			c.OldList.MoveToFront(ent)
			return false
		}
	}

	ent := &entry{key, value, time.Now(), false}
	e := c.OldList.PushFront(ent)
	c.items[key] = e

	evicted := c.YoungList.Len()+c.OldList.Len() > c.size
	if evicted {
		c.removeOldest()
	}

	return evicted
}

func (c *YoungOldLRU) Get(key any) (value any, ok bool) {
	if ent, ok := c.items[key]; ok {
		kv := ent.Value.(*entry)
		if kv == nil {
			return nil, false
		}
		if kv.flag {
			c.YoungList.MoveToFront(ent)
			return kv.value, true
		} else {
			if kv.addAt.Before(time.Now().Add(-c.stayTime)) {
				// 如果在 old 队列中的时间超过了 stayTime，允许晋升到 young 队列
				c.promote(key)
				return kv.value, true
			}
			// 否则只移动到 old 队列头部
			c.OldList.MoveToFront(ent)
			return kv.value, true
		}
	}
	return nil, false
}

func (c *YoungOldLRU) Contains(key any) (ok bool) {
	_, ok = c.items[key]
	return ok
}

func (c *YoungOldLRU) Peek(key any) (value any, ok bool) {
	var ent *list.Element
	if ent, ok = c.items[key]; ok {
		return ent.Value.(*entry).value, true
	}
	return nil, ok
}

// 修复Remove方法
func (c *YoungOldLRU) Remove(key any) (present bool) {
	if ent, ok := c.items[key]; ok {
		kv := ent.Value.(*entry)
		if kv == nil {
			return false
		}
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
	return false
}

// 修复Resize方法
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
		ent := c.YoungList.Back()
		kv := ent.Value.(*entry)
		if kv != nil {
			kv.flag = false
			kv.addAt = time.Now()
			newent := c.OldList.PushFront(kv)
			c.items[kv.key] = newent
		}
		c.YoungList.Remove(ent)
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

func (c *YoungOldLRU) RemoveOldest() (key, value any, ok bool) {
	ent := c.OldList.Back()
	if ent != nil {
		c.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

func (c *YoungOldLRU) GetOldest() (key, value any, ok bool) {
	ent := c.OldList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

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

func (c *YoungOldLRU) Len() int {
	return len(c.items)
}

func (c *YoungOldLRU) promote(key any) {
	ent := c.items[key]
	if ent != nil {
		kv := ent.Value.(*entry)
		kv.flag = true
		newent := c.YoungList.PushFront(kv)
		c.items[key] = newent
		c.OldList.Remove(ent)
		if c.YoungList.Len() > c.youngListSize {
			ent := c.YoungList.Back()
			kv := ent.Value.(*entry)
			kv.flag = false
			kv.addAt = time.Now()
			c.YoungList.Remove(ent)
			newent := c.OldList.PushFront(kv)
			c.items[kv.key] = newent
		}
	}
}

func (c *YoungOldLRU) removeOldest() {
	ent := c.OldList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

func (c *YoungOldLRU) removeElement(e *list.Element) {
	c.OldList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}
