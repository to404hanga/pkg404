package lru

import (
	"sync"
	"time"

	"github.com/to404hanga/pkg404/cachex/lru/internal/interfaces"
	"github.com/to404hanga/pkg404/cachex/lru/internal/simple_lru"
	"github.com/to404hanga/pkg404/cachex/lru/internal/young_old_lru"
)

const (
	DEFAULT_EVICTED_BUFFER_SIZE = 16
)

type Cache struct {
	lru         interfaces.LRUCache
	evictedKeys []any
	evictedVals []any
	onEvictedCB func(k, v any)
	lock        sync.RWMutex
}

func NewYoungOldLRU(size, youngSize int, stayTime time.Duration) (*Cache, error) {
	return NewYoungOldLRUWithEvict(size, youngSize, stayTime, nil)
}

func NewYoungOldLRUWithEvict(size, youngSize int, stayTime time.Duration, onEvicted func(k, v any)) (c *Cache, err error) {
	c = &Cache{
		onEvictedCB: onEvicted,
	}
	if onEvicted != nil {
		c.initEvictBuffers()
		onEvicted = c.onEvicted
	}
	c.lru, err = young_old_lru.NewYoungOldLRU(size, youngSize, stayTime, onEvicted)

	return
}

func NewSimpleLRU(size int) (*Cache, error) {
	return NewSimpleLRUWithEvict(size, nil)
}

func NewSimpleLRUWithEvict(size int, onEvicted func(k, v any)) (c *Cache, err error) {
	c = &Cache{
		onEvictedCB: onEvicted,
	}
	if onEvicted != nil {
		c.initEvictBuffers()
		onEvicted = c.onEvicted
	}
	c.lru, err = simple_lru.NewLRU(size, onEvicted)
	return
}

func (c *Cache) initEvictBuffers() {
	c.evictedKeys = make([]any, 0, DEFAULT_EVICTED_BUFFER_SIZE)
	c.evictedVals = make([]any, 0, DEFAULT_EVICTED_BUFFER_SIZE)
}

func (c *Cache) onEvicted(k, v any) {
	c.evictedKeys = append(c.evictedKeys, k)
	c.evictedVals = append(c.evictedVals, v)
}

func (c *Cache) Purge() {
	var ks, vs []any
	c.lock.Lock()
	c.lru.Purge()
	if c.onEvictedCB != nil && len(c.evictedKeys) > 0 {
		ks, vs = c.evictedKeys, c.evictedVals
		c.initEvictBuffers()
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil {
		for i := 0; i < len(ks); i++ {
			c.onEvictedCB(ks[i], vs[i])
		}
	}
}

func (c *Cache) Add(key, value any) (evicted bool) {
	var k, v any
	c.lock.Lock()
	evicted = c.lru.Add(key, value)
	if c.onEvictedCB != nil && evicted {
		k, v = c.evictedKeys[0], c.evictedVals[0]
		c.evictedKeys, c.evictedVals = c.evictedKeys[1:], c.evictedVals[1:]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && evicted {
		c.onEvictedCB(k, v)
	}
	return
}

func (c *Cache) Get(key any) (value any, ok bool) {
	c.lock.Lock()
	value, ok = c.lru.Get(key)
	c.lock.Unlock()
	return value, ok
}

func (c *Cache) Contains(key any) bool {
	c.lock.RLock()
	containKey := c.lru.Contains(key)
	c.lock.RUnlock()
	return containKey
}

func (c *Cache) Peek(key any) (value any, ok bool) {
	c.lock.RLock()
	value, ok = c.lru.Peek(key)
	c.lock.RUnlock()
	return value, ok
}

func (c *Cache) ContainsOrAdd(key, value any) (ok, evicted bool) {
	var k, v any
	c.lock.Lock()
	if c.lru.Contains(key) {
		c.lock.Unlock()
		return true, false
	}
	evicted = c.lru.Add(key, value)
	if c.onEvictedCB != nil && evicted {
		k, v = c.evictedKeys[0], c.evictedVals[0]
		c.evictedKeys, c.evictedVals = c.evictedKeys[:0], c.evictedVals[:0]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && evicted {
		c.onEvictedCB(k, v)
	}
	return false, evicted
}

func (c *Cache) PeekOrAdd(key, value any) (previous any, ok, evicted bool) {
	var k, v any
	c.lock.Lock()
	previous, ok = c.lru.Peek(key)
	if ok {
		c.lock.Unlock()
		return previous, true, false
	}
	evicted = c.lru.Add(key, value)
	if c.onEvictedCB != nil && evicted {
		k, v = c.evictedKeys[0], c.evictedVals[0]
		c.evictedKeys, c.evictedVals = c.evictedKeys[:0], c.evictedVals[:0]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && evicted {
		c.onEvictedCB(k, v)
	}
	return nil, false, evicted
}

func (c *Cache) Remove(key any) (present bool) {
	var k, v any
	c.lock.Lock()
	present = c.lru.Remove(key)
	if c.onEvictedCB != nil && present {
		k, v = c.evictedKeys[0], c.evictedVals[0]
		c.evictedKeys, c.evictedVals = c.evictedKeys[:0], c.evictedVals[:0]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && present {
		c.onEvictedCB(k, v)
	}
	return
}

func (c *Cache) Resize(opts ...*interfaces.SizeOptions) (evicted int) {
	var ks, vs []any
	c.lock.Lock()
	evicted = c.lru.Resize(opts...)
	if c.onEvictedCB != nil && evicted > 0 {
		ks, vs = c.evictedKeys, c.evictedVals
		c.initEvictBuffers()
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && evicted > 0 {
		for i := 0; i < len(ks); i++ {
			c.onEvictedCB(ks[i], vs[i])
		}
	}
	return evicted
}

func (c *Cache) RemoveOldest() (key, value any, ok bool) {
	var k, v any
	c.lock.Lock()
	key, value, ok = c.lru.RemoveOldest()
	if c.onEvictedCB != nil && ok {
		k, v = c.evictedKeys[0], c.evictedVals[0]
		c.evictedKeys, c.evictedVals = c.evictedKeys[:0], c.evictedVals[:0]
	}
	c.lock.Unlock()
	if c.onEvictedCB != nil && ok {
		c.onEvictedCB(k, v)
	}
	return
}

func (c *Cache) GetOldest() (key, value any, ok bool) {
	c.lock.RLock()
	key, value, ok = c.lru.GetOldest()
	c.lock.RUnlock()
	return
}

func (c *Cache) Keys() []any {
	c.lock.RLock()
	keys := c.lru.Keys()
	c.lock.RUnlock()
	return keys
}

func (c *Cache) Len() int {
	c.lock.RLock()
	length := c.lru.Len()
	c.lock.RUnlock()
	return length
}
