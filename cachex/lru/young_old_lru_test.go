package lru

import (
	"testing"
	"time"

	"github.com/to404hanga/pkg404/cachex/lru/internal/interfaces"
	"github.com/to404hanga/pkg404/cachex/lru/internal/young_old_lru"
)

// BenchmarkYoungOldLRU_Rand_YoungOldLRU 随机访问模式的基准测试
func BenchmarkYoungOldLRU_Rand_YoungOldLRU(b *testing.B) {
	l, err := NewYoungOldLRU(8192, 2048, 100*time.Millisecond)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.Add(trace[i], trace[i])
		} else {
			_, ok := l.Get(trace[i])
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

// BenchmarkYoungOldLRU_Freq_YoungOldLRU 频率访问模式的基准测试
func BenchmarkYoungOldLRU_Freq_YoungOldLRU(b *testing.B) {
	l, err := NewYoungOldLRU(8192, 2048, 100*time.Millisecond)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = getRand(b) % 16384
		} else {
			trace[i] = getRand(b) % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Add(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		_, ok := l.Get(trace[i])
		if ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

// TestYoungOldLRU 主要功能测试
func TestYoungOldLRU(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k any, v any) {
		if k != v {
			t.Fatalf("Evict values not equal (%v!=%v)", k, v)
		}
		evictCounter++
	}
	l, err := NewYoungOldLRUWithEvict(128, 32, 50*time.Millisecond, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// 添加256个元素，应该触发淘汰
	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}

	if evictCounter != 128 {
		t.Fatalf("bad evict count: %v", evictCounter)
	}

	// 验证keys的顺序和值
	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || v != k {
			t.Fatalf("bad key: %v", k)
		}
		// 重置计数器位置，因为Get可能改变顺序
		_ = i
	}

	// 验证前128个元素被淘汰
	for i := 0; i < 128; i++ {
		_, ok := l.Get(i)
		if ok {
			t.Fatalf("should be evicted")
		}
	}

	// 验证后128个元素仍存在
	for i := 128; i < 256; i++ {
		_, ok := l.Get(i)
		if !ok {
			t.Fatalf("should not be evicted")
		}
	}

	// 删除一些元素
	for i := 128; i < 192; i++ {
		l.Remove(i)
		_, ok := l.Get(i)
		if ok {
			t.Fatalf("should be deleted")
		}
	}

	// 清空缓存
	l.Purge()
	if l.Len() != 0 {
		t.Fatalf("bad len: %v", l.Len())
	}
	if _, ok := l.Get(200); ok {
		t.Fatalf("should contain nothing")
	}
}

// TestYoungOldLRUAdd 测试Add方法
func TestYoungOldLRUAdd(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k any, v any) {
		evictCounter++
	}

	l, err := NewYoungOldLRUWithEvict(1, 1, 50*time.Millisecond, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if l.Add(1, 1) == true || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	if l.Add(2, 2) == false || evictCounter != 1 {
		t.Errorf("should have an eviction")
	}
}

// TestYoungOldLRUContains 测试Contains方法
func TestYoungOldLRUContains(t *testing.T) {
	l, err := NewYoungOldLRU(2, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add(1, 1)
	l.Add(2, 2)
	if !l.Contains(1) {
		t.Errorf("1 should be contained")
	}

	l.Add(3, 3)
	if l.Contains(1) {
		t.Errorf("Contains should not have updated recent-ness of 1")
	}
}

// TestYoungOldLRUPeek 测试Peek方法
func TestYoungOldLRUPeek(t *testing.T) {
	l, err := NewYoungOldLRU(2, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add(1, 1)
	l.Add(2, 2)
	if v, ok := l.Peek(1); !ok || v != 1 {
		t.Errorf("1 should be set to 1: %v, %v", v, ok)
	}

	l.Add(3, 3)
	if l.Contains(1) {
		t.Errorf("should not have updated recent-ness of 1")
	}
}

// TestYoungOldLRUResize 测试Resize方法
func TestYoungOldLRUResize(t *testing.T) {
	onEvictCounter := 0
	onEvicted := func(k any, v any) {
		onEvictCounter++
	}
	l, err := NewYoungOldLRUWithEvict(4, 2, 50*time.Millisecond, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add(1, 1)
	l.Add(2, 2)
	l.Add(3, 3)
	l.Add(4, 4)

	// 缩小缓存大小
	evicted := l.Resize(interfaces.WithSize(2), interfaces.WithYoungListSize(1))
	if evicted != 2 {
		t.Errorf("2 elements should have been evicted: %v", evicted)
	}
	if onEvictCounter != 2 {
		t.Errorf("onEvicted should have been called 2 times: %v", onEvictCounter)
	}

	if l.Len() != 2 {
		t.Errorf("Cache should have 2 elements: %v", l.Len())
	}

	// 扩大缓存大小
	evicted = l.Resize(interfaces.WithSize(4), interfaces.WithYoungListSize(2))
	if evicted != 0 {
		t.Errorf("0 elements should have been evicted: %v", evicted)
	}

	l.Add(5, 5)
	l.Add(6, 6)
	if l.Len() != 4 {
		t.Errorf("Cache should have 4 elements: %v", l.Len())
	}
}

// TestYoungOldLRURemove 测试Remove方法
func TestYoungOldLRURemove(t *testing.T) {
	l, err := NewYoungOldLRU(2, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add(1, 1)
	l.Add(2, 2)

	// 移除存在的元素
	if !l.Remove(1) {
		t.Errorf("should have removed 1")
	}
	if l.Contains(1) {
		t.Errorf("should not contain 1")
	}

	// 移除不存在的元素
	if l.Remove(3) {
		t.Errorf("should not have removed 3")
	}
}

// TestYoungOldLRURemoveOldest 测试RemoveOldest方法
func TestYoungOldLRURemoveOldest(t *testing.T) {
	l, err := NewYoungOldLRU(2, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// 空缓存测试
	key, value, ok := l.RemoveOldest()
	if ok {
		t.Errorf("should not have removed from empty cache")
	}

	l.Add(1, 1)
	l.Add(2, 2)

	// 移除最老的元素
	key, value, ok = l.RemoveOldest()
	if !ok {
		t.Errorf("should have removed oldest")
	}
	if key != 1 || value != 1 {
		t.Errorf("should have removed 1,1 got %v,%v", key, value)
	}
}

// TestYoungOldLRUGetOldest 测试GetOldest方法
func TestYoungOldLRUGetOldest(t *testing.T) {
	l, err := NewYoungOldLRU(2, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// 空缓存测试
	key, value, ok := l.GetOldest()
	if ok {
		t.Errorf("should not have gotten from empty cache")
	}

	l.Add(1, 1)
	l.Add(2, 2)

	// 获取最老的元素
	key, value, ok = l.GetOldest()
	if !ok {
		t.Errorf("should have gotten oldest")
	}
	if key != 1 || value != 1 {
		t.Errorf("should have gotten 1,1 got %v,%v", key, value)
	}

	// 验证元素仍在缓存中
	if !l.Contains(1) {
		t.Errorf("should still contain 1")
	}
}

// TestYoungOldLRUKeys 测试Keys方法
func TestYoungOldLRUKeys(t *testing.T) {
	l, err := NewYoungOldLRU(3, 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// 空缓存测试
	keys := l.Keys()
	if len(keys) != 0 {
		t.Errorf("should have no keys")
	}

	l.Add(1, 1)
	l.Add(2, 2)
	l.Add(3, 3)

	keys = l.Keys()
	if len(keys) != 3 {
		t.Errorf("should have 3 keys")
	}
}

// TestYoungOldPromotion 测试Young-Old晋升机制
func TestYoungOldPromotion(t *testing.T) {
	l, err := NewYoungOldLRU(5, 2, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// 添加元素，应该都在old队列
	l.Add(1, 1)
	l.Add(2, 2)
	l.Add(3, 3)

	if l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len() != 0 {
		t.Errorf("young list should be empty, got %d", l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len())
	}
	if l.lru.(*young_old_lru.YoungOldLRU).OldList.Len() != 3 {
		t.Errorf("old list should have 3 elements, got %d", l.lru.(*young_old_lru.YoungOldLRU).OldList.Len())
	}

	// 等待超过stayTime
	time.Sleep(60 * time.Millisecond)

	// 访问元素1，应该晋升到young队列
	l.Get(1)

	if l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len() != 1 {
		t.Errorf("young list should have 1 element after promotion, got %d", l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len())
	}
	if l.lru.(*young_old_lru.YoungOldLRU).OldList.Len() != 2 {
		t.Errorf("old list should have 2 elements after promotion, got %d", l.lru.(*young_old_lru.YoungOldLRU).OldList.Len())
	}

	// 再次访问元素1，应该仍在young队列
	l.Get(1)

	if l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len() != 1 {
		t.Errorf("young list should still have 1 element, got %d", l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len())
	}
}

// TestYoungOldStayTime 测试stayTime机制
func TestYoungOldStayTime(t *testing.T) {
	l, err := NewYoungOldLRU(5, 2, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// 添加元素
	l.Add(1, 1)

	// 立即访问，不应该晋升
	l.Get(1)
	if l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len() != 0 {
		t.Errorf("should not promote before stayTime")
	}

	// 等待超过stayTime
	time.Sleep(110 * time.Millisecond)

	// 现在访问应该晋升
	l.Get(1)
	if l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len() != 1 {
		t.Errorf("should promote after stayTime")
	}
}

// TestYoungListOverflow 测试young队列溢出处理
func TestYoungListOverflow(t *testing.T) {
	l, err := NewYoungOldLRU(5, 2, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// 添加元素并等待
	for i := 1; i <= 3; i++ {
		l.Add(i, i)
	}
	time.Sleep(60 * time.Millisecond)

	// 晋升3个元素到young队列，应该触发溢出处理
	for i := 1; i <= 3; i++ {
		l.Get(i)
	}

	// young队列应该只有2个元素（youngListSize限制）
	if l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len() != 2 {
		t.Errorf("young list should have 2 elements, got %d", l.lru.(*young_old_lru.YoungOldLRU).YoungList.Len())
	}
	// old队列应该有1个元素（被从young队列降级的）
	if l.lru.(*young_old_lru.YoungOldLRU).OldList.Len() != 1 {
		t.Errorf("old list should have 1 element, got %d", l.lru.(*young_old_lru.YoungOldLRU).OldList.Len())
	}
}

// TestNewYoungOldLRUErrors 测试构造函数错误处理
func TestNewYoungOldLRUErrors(t *testing.T) {
	// 测试size <= 0
	_, err := NewYoungOldLRU(0, 1, time.Second)
	if err == nil {
		t.Error("should return error for size <= 0")
	}

	// 测试youngListSize <= 0
	_, err = NewYoungOldLRU(10, 0, time.Second)
	if err == nil {
		t.Error("should return error for youngListSize <= 0")
	}

	// 测试youngListSize > size
	_, err = NewYoungOldLRU(10, 15, time.Second)
	if err == nil {
		t.Error("should return error for youngListSize > size")
	}

	// 测试正常情况
	_, err = NewYoungOldLRU(10, 3, time.Second)
	if err != nil {
		t.Errorf("should not return error for valid parameters: %v", err)
	}
}
