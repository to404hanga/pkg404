package simple_lru

import (
	"container/list"
	"errors"

	"github.com/to404hanga/pkg404/cachex/lru/internal/interfaces"
)

type EvictCallback func(key any, value any)

type LRU struct {
	size      int
	evictList *list.List
	items     map[any]*list.Element
	onEvict   EvictCallback
}

type entry struct {
	key   any
	value any
}

func NewLRU(size int, onEvict EvictCallback) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}
	c := &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[any]*list.Element),
		onEvict:   onEvict,
	}
	return c, nil
}

func (c *LRU) Purge() {
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entry).value)
		}
		delete(c.items, k)
	}
	c.evictList.Init()
}

func (c *LRU) Add(key, value any) (evicted bool) {

	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		return false
	}

	ent := &entry{key, value}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry

	evict := c.evictList.Len() > c.size

	if evict {
		c.removeOldest()
	}
	return evict
}

func (c *LRU) Get(key any) (value any, ok bool) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		if ent.Value.(*entry) == nil {
			return nil, false
		}
		return ent.Value.(*entry).value, true
	}
	return
}

func (c *LRU) Contains(key any) (ok bool) {
	_, ok = c.items[key]
	return ok
}

func (c *LRU) Peek(key any) (value any, ok bool) {
	var ent *list.Element
	if ent, ok = c.items[key]; ok {
		return ent.Value.(*entry).value, true
	}
	return nil, ok
}

func (c *LRU) Remove(key any) (present bool) {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

func (c *LRU) RemoveOldest() (key, value any, ok bool) {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

func (c *LRU) GetOldest() (key, value any, ok bool) {
	ent := c.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

func (c *LRU) Keys() []any {
	keys := make([]any, len(c.items))
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

func (c *LRU) Len() int {
	return c.evictList.Len()
}

func (c *LRU) Resize(opts ...*interfaces.SizeOptions) (evicted int) {
	var size int
	for _, opt := range opts {
		if opt.Key == "size" {
			size = opt.Value
		}
	}
	if size <= 0 {
		c.Purge()
		return -1
	}

	diff := c.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.removeOldest()
	}
	c.size = size
	return diff
}

func (c *LRU) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

func (c *LRU) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}
