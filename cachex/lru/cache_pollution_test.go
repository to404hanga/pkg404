package lru

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkCachePollution_ScanAttack_SimpleLRU 测试SimpleLRU在扫描攻击场景下的表现
func BenchmarkCachePollution_ScanAttack_SimpleLRU(b *testing.B) {
	cache, err := NewSimpleLRU(1000)
	if err != nil {
		b.Fatal(err)
	}

	// 预热缓存：填入热点数据
	for i := 0; i < 100; i++ {
		cache.Add(fmt.Sprintf("hot_%d", i), i)
	}

	hitCount := 0
	missCount := 0

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if i%10 == 0 {
			// 10%的时间访问热点数据
			key := fmt.Sprintf("hot_%d", i%100)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
			}
		} else {
			// 90%的时间进行扫描攻击（访问大量冷数据）
			key := fmt.Sprintf("scan_%d", i)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
				cache.Add(key, i)
			}
		}
	}

	b.StopTimer()
	hitRate := float64(hitCount) / float64(hitCount+missCount) * 100
	b.Logf("Hit rate: %.2f%% (hits: %d, misses: %d)", hitRate, hitCount, missCount)
}

// BenchmarkCachePollution_ScanAttack_YoungOldLRU 测试YoungOldLRU在扫描攻击场景下的表现
func BenchmarkCachePollution_ScanAttack_YoungOldLRU(b *testing.B) {
	cache, err := NewYoungOldLRU(1000, 500, 10*time.Millisecond)
	if err != nil {
		b.Fatal(err)
	}

	// 预热缓存：填入热点数据
	for i := 0; i < 100; i++ {
		cache.Add(fmt.Sprintf("hot_%d", i), i)
	}

	hitCount := 0
	missCount := 0

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if i%10 == 0 {
			// 10%的时间访问热点数据
			key := fmt.Sprintf("hot_%d", i%100)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
			}
		} else {
			// 90%的时间进行扫描攻击
			key := fmt.Sprintf("scan_%d", i)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
				cache.Add(key, i)
			}
		}
	}

	b.StopTimer()
	hitRate := float64(hitCount) / float64(hitCount+missCount) * 100
	b.Logf("Hit rate: %.2f%% (hits: %d, misses: %d)", hitRate, hitCount, missCount)
}

// BenchmarkCachePollution_BurstTraffic_SimpleLRU 测试SimpleLRU在突发流量场景下的表现
func BenchmarkCachePollution_BurstTraffic_SimpleLRU(b *testing.B) {
	cache, err := NewSimpleLRU(1000)
	if err != nil {
		b.Fatal(err)
	}

	// 预热缓存：填入稳定的热点数据
	for i := 0; i < 200; i++ {
		cache.Add(fmt.Sprintf("stable_%d", i), i)
	}

	hitCount := 0
	missCount := 0

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if i < b.N/3 {
			// 前1/3时间：正常访问模式
			key := fmt.Sprintf("stable_%d", i%200)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
			}
		} else if i < 2*b.N/3 {
			// 中间1/3时间：突发流量（大量新数据）
			key := fmt.Sprintf("burst_%d", i)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
				cache.Add(key, i)
			}
		} else {
			// 后1/3时间：回到正常访问模式
			key := fmt.Sprintf("stable_%d", i%200)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
			}
		}
	}

	b.StopTimer()
	hitRate := float64(hitCount) / float64(hitCount+missCount) * 100
	b.Logf("Hit rate: %.2f%% (hits: %d, misses: %d)", hitRate, hitCount, missCount)
}

// BenchmarkCachePollution_BurstTraffic_YoungOldLRU 测试YoungOldLRU在突发流量场景下的表现
func BenchmarkCachePollution_BurstTraffic_YoungOldLRU(b *testing.B) {
	cache, err := NewYoungOldLRU(1000, 200, 10*time.Millisecond)
	if err != nil {
		b.Fatal(err)
	}

	// 预热缓存：填入稳定的热点数据
	for i := 0; i < 200; i++ {
		cache.Add(fmt.Sprintf("stable_%d", i), i)
	}

	hitCount := 0
	missCount := 0

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if i < b.N/3 {
			// 前1/3时间：正常访问模式
			key := fmt.Sprintf("stable_%d", i%200)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
			}
		} else if i < 2*b.N/3 {
			// 中间1/3时间：突发流量
			key := fmt.Sprintf("burst_%d", i)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
				cache.Add(key, i)
			}
		} else {
			// 后1/3时间：回到正常访问模式
			key := fmt.Sprintf("stable_%d", i%200)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
			}
		}
	}

	b.StopTimer()
	hitRate := float64(hitCount) / float64(hitCount+missCount) * 100
	b.Logf("Hit rate: %.2f%% (hits: %d, misses: %d)", hitRate, hitCount, missCount)
}

// BenchmarkCachePollution_WorkingSetShift_SimpleLRU 测试SimpleLRU在工作集变化场景下的表现
func BenchmarkCachePollution_WorkingSetShift_SimpleLRU(b *testing.B) {
	cache, err := NewSimpleLRU(1000)
	if err != nil {
		b.Fatal(err)
	}

	// 预热缓存：填入第一个工作集
	for i := 0; i < 300; i++ {
		cache.Add(fmt.Sprintf("workset1_%d", i), i)
	}

	hitCount := 0
	missCount := 0

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if i < b.N/2 {
			// 前半段：访问第一个工作集
			key := fmt.Sprintf("workset1_%d", i%300)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
			}
		} else {
			// 后半段：工作集完全切换到第二个工作集
			key := fmt.Sprintf("workset2_%d", i%300)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
				cache.Add(key, i)
			}
		}
	}

	b.StopTimer()
	hitRate := float64(hitCount) / float64(hitCount+missCount) * 100
	b.Logf("Hit rate: %.2f%% (hits: %d, misses: %d)", hitRate, hitCount, missCount)
}

// BenchmarkCachePollution_WorkingSetShift_YoungOldLRU 测试YoungOldLRU在工作集变化场景下的表现
func BenchmarkCachePollution_WorkingSetShift_YoungOldLRU(b *testing.B) {
	cache, err := NewYoungOldLRU(1000, 200, 10*time.Millisecond)
	if err != nil {
		b.Fatal(err)
	}

	// 预热缓存：填入第一个工作集
	for i := 0; i < 300; i++ {
		cache.Add(fmt.Sprintf("workset1_%d", i), i)
	}

	hitCount := 0
	missCount := 0

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if i < b.N/2 {
			// 前半段：访问第一个工作集
			key := fmt.Sprintf("workset1_%d", i%300)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
			}
		} else {
			// 后半段：工作集完全切换
			key := fmt.Sprintf("workset2_%d", i%300)
			if _, ok := cache.Get(key); ok {
				hitCount++
			} else {
				missCount++
				cache.Add(key, i)
			}
		}
	}

	b.StopTimer()
	hitRate := float64(hitCount) / float64(hitCount+missCount) * 100
	b.Logf("Hit rate: %.2f%% (hits: %d, misses: %d)", hitRate, hitCount, missCount)
}

// TestRunCachePollutionComparison 运行缓存污染场景的性能比较测试
func TestRunCachePollutionComparison(t *testing.T) {
	t.Log("=== 缓存污染场景性能测试 ===")
	t.Log("")
	t.Log("测试场景说明：")
	t.Log("1. 扫描攻击：90%访问冷数据，10%访问热点数据")
	t.Log("2. 突发流量：正常访问 -> 突发新数据 -> 恢复正常访问")
	t.Log("3. 工作集变化：第一个工作集 -> 完全切换到第二个工作集")
	t.Log("")
	t.Log("运行命令：")
	t.Log("go test -bench=BenchmarkCachePollution -benchmem -v")
	t.Log("")
	t.Log("预期结果：")
	t.Log("- YoungOldLRU在扫描攻击场景下应该表现更好（更好的污染抵抗能力）")
	t.Log("- YoungOldLRU在突发流量场景下应该能更好地保护热点数据")
	t.Log("- 工作集变化场景下两者表现可能相近，但YoungOldLRU适应性更强")
}
