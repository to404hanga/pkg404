package young_old_lru

import (
	"testing"
	"time"

	"github.com/to404hanga/pkg404/cachex/lru/internal/interfaces"
)

func TestNewYoungOldLRU(t *testing.T) {
	// 测试正常创建
	lru, err := NewYoungOldLRU(10, 3, time.Second, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if lru.size != 10 {
		t.Errorf("Expected size 10, got %d", lru.size)
	}
	if lru.youngListSize != 3 {
		t.Errorf("Expected youngListSize 3, got %d", lru.youngListSize)
	}

	// 测试无效参数
	_, err = NewYoungOldLRU(0, 3, time.Second, nil)
	if err == nil {
		t.Error("Expected error for size <= 0")
	}

	_, err = NewYoungOldLRU(10, 0, time.Second, nil)
	if err == nil {
		t.Error("Expected error for youngListSize <= 0")
	}

	_, err = NewYoungOldLRU(10, 15, time.Second, nil)
	if err == nil {
		t.Error("Expected error for youngListSize > size")
	}
}

func TestAdd(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, 100*time.Millisecond, nil)

	// 测试添加新元素
	evicted := lru.Add("key1", "value1")
	if evicted {
		t.Error("Expected no eviction for first add")
	}
	if lru.Len() != 1 {
		t.Errorf("Expected length 1, got %d", lru.Len())
	}

	// 测试添加重复key
	evicted = lru.Add("key1", "newvalue1")
	if evicted {
		t.Error("Expected no eviction for duplicate key")
	}
	if lru.Len() != 1 {
		t.Errorf("Expected length 1, got %d", lru.Len())
	}

	// 测试填满缓存
	for i := 2; i <= 5; i++ {
		lru.Add(i, i*10)
	}
	if lru.Len() != 5 {
		t.Errorf("Expected length 5, got %d", lru.Len())
	}

	// 测试触发淘汰
	evicted = lru.Add("key6", "value6")
	if !evicted {
		t.Error("Expected eviction when cache is full")
	}
	if lru.Len() != 5 {
		t.Errorf("Expected length 5, got %d", lru.Len())
	}
}

func TestGet(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, 100*time.Millisecond, nil)

	// 测试获取不存在的key
	value, ok := lru.Get("nonexistent")
	if ok {
		t.Error("Expected false for nonexistent key")
	}
	if value != nil {
		t.Error("Expected nil value for nonexistent key")
	}

	// 添加元素并测试获取
	lru.Add("key1", "value1")
	value, ok = lru.Get("key1")
	if !ok {
		t.Error("Expected true for existing key")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}

func TestPromote(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, 50*time.Millisecond, nil)

	// 添加元素到old队列
	lru.Add("key1", "value1")

	// 等待超过stayTime
	time.Sleep(60 * time.Millisecond)

	// 访问应该触发晋升
	value, ok := lru.Get("key1")
	if !ok {
		t.Error("Expected to find key1")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}

	// 验证元素已经在young队列中
	if lru.YoungList.Len() != 1 {
		t.Errorf("Expected 1 element in young list, got %d", lru.YoungList.Len())
	}
}

func TestRemove(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, time.Second, nil)

	// 测试移除不存在的key
	present := lru.Remove("nonexistent")
	if present {
		t.Error("Expected false for nonexistent key")
	}

	// 添加元素并移除
	lru.Add("key1", "value1")
	present = lru.Remove("key1")
	if !present {
		t.Error("Expected true for existing key")
	}
	if lru.Len() != 0 {
		t.Errorf("Expected length 0 after removal, got %d", lru.Len())
	}
}

func TestEvictCallback(t *testing.T) {
	evictedKeys := make([]any, 0)
	evictedValues := make([]any, 0)

	onEvict := func(key, value any) {
		evictedKeys = append(evictedKeys, key)
		evictedValues = append(evictedValues, value)
	}

	lru, _ := NewYoungOldLRU(2, 1, time.Second, onEvict)

	// 填满缓存
	lru.Add("key1", "value1")
	lru.Add("key2", "value2")

	// 添加第三个元素应该触发淘汰
	lru.Add("key3", "value3")

	if len(evictedKeys) != 1 {
		t.Errorf("Expected 1 evicted key, got %d", len(evictedKeys))
	}
	if evictedKeys[0] != "key1" {
		t.Errorf("Expected evicted key 'key1', got %v", evictedKeys[0])
	}
	if evictedValues[0] != "value1" {
		t.Errorf("Expected evicted value 'value1', got %v", evictedValues[0])
	}
}

func TestPeek(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, time.Second, nil)

	// 测试peek不存在的key
	value, ok := lru.Peek("nonexistent")
	if ok {
		t.Error("Expected false for nonexistent key")
	}

	// 添加元素并peek
	lru.Add("key1", "value1")
	value, ok = lru.Peek("key1")
	if !ok {
		t.Error("Expected true for existing key")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}

func TestContains(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, time.Second, nil)

	// 测试不存在的key
	if lru.Contains("nonexistent") {
		t.Error("Expected false for nonexistent key")
	}

	// 添加元素并测试
	lru.Add("key1", "value1")
	if !lru.Contains("key1") {
		t.Error("Expected true for existing key")
	}
}

func TestRemoveOldest(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, time.Second, nil)

	// 测试空缓存
	key, value, ok := lru.RemoveOldest()
	if ok {
		t.Error("Expected false for empty cache")
	}

	// 添加元素并移除最老的
	lru.Add("key1", "value1")
	lru.Add("key2", "value2")

	key, value, ok = lru.RemoveOldest()
	if !ok {
		t.Error("Expected true for non-empty cache")
	}
	if key != "key1" {
		t.Errorf("Expected oldest key 'key1', got %v", key)
	}
	if value != "value1" {
		t.Errorf("Expected oldest value 'value1', got %v", value)
	}
}

func TestGetOldest(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, time.Second, nil)

	// 测试空缓存
	key, value, ok := lru.GetOldest()
	if ok {
		t.Error("Expected false for empty cache")
	}

	// 添加元素并获取最老的
	lru.Add("key1", "value1")
	lru.Add("key2", "value2")

	key, value, ok = lru.GetOldest()
	if !ok {
		t.Error("Expected true for non-empty cache")
	}
	if key != "key1" {
		t.Errorf("Expected oldest key 'key1', got %v", key)
	}
	if value != "value1" {
		t.Errorf("Expected oldest value 'value1', got %v", value)
	}

	// 验证元素仍在缓存中
	if lru.Len() != 2 {
		t.Errorf("Expected length 2 after GetOldest, got %d", lru.Len())
	}
}

func TestKeys(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, time.Second, nil)

	// 测试空缓存
	keys := lru.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected 0 keys for empty cache, got %d", len(keys))
	}

	// 添加元素并获取keys
	lru.Add("key1", "value1")
	lru.Add("key2", "value2")
	lru.Add("key3", "value3")

	keys = lru.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
}

func TestPurge(t *testing.T) {
	evictedCount := 0
	onEvict := func(key, value any) {
		evictedCount++
	}

	lru, _ := NewYoungOldLRU(5, 2, time.Second, onEvict)

	// 添加一些元素
	lru.Add("key1", "value1")
	lru.Add("key2", "value2")
	lru.Add("key3", "value3")

	// 清空缓存
	lru.Purge()

	if lru.Len() != 0 {
		t.Errorf("Expected length 0 after purge, got %d", lru.Len())
	}
	if evictedCount != 3 {
		t.Errorf("Expected 3 evictions, got %d", evictedCount)
	}
}

func TestResize(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, time.Second, nil)

	// 添加元素
	for i := 1; i <= 5; i++ {
		lru.Add(i, i*10)
	}

	// 缩小缓存
	evicted := lru.Resize(interfaces.WithSize(3), interfaces.WithYoungListSize(1))
	if evicted != 2 {
		t.Errorf("Expected 2 evictions, got %d", evicted)
	}
	if lru.Len() != 3 {
		t.Errorf("Expected length 3 after resize, got %d", lru.Len())
	}
	if lru.size != 3 {
		t.Errorf("Expected size 3, got %d", lru.size)
	}
	if lru.youngListSize != 1 {
		t.Errorf("Expected youngListSize 1, got %d", lru.youngListSize)
	}
}

func TestYoungOldBehavior(t *testing.T) {
	lru, _ := NewYoungOldLRU(5, 2, 50*time.Millisecond, nil)

	// 添加元素，应该都在old队列
	lru.Add("key1", "value1")
	lru.Add("key2", "value2")
	lru.Add("key3", "value3")

	if lru.OldList.Len() != 3 {
		t.Errorf("Expected 3 elements in old list, got %d", lru.OldList.Len())
	}
	if lru.YoungList.Len() != 0 {
		t.Errorf("Expected 0 elements in young list, got %d", lru.YoungList.Len())
	}

	// 等待超过stayTime
	time.Sleep(60 * time.Millisecond)

	// 访问key1，应该晋升到young队列
	lru.Get("key1")

	if lru.YoungList.Len() != 1 {
		t.Errorf("Expected 1 element in young list after promotion, got %d", lru.YoungList.Len())
	}
	if lru.OldList.Len() != 2 {
		t.Errorf("Expected 2 elements in old list after promotion, got %d", lru.OldList.Len())
	}

	// 再次访问key1，应该仍在young队列
	lru.Get("key1")

	if lru.YoungList.Len() != 1 {
		t.Errorf("Expected 1 element in young list, got %d", lru.YoungList.Len())
	}
}
