package lru

import (
	"fmt"
	"testing"
	"time"
)

// 性能比较测试 - 随机访问模式
func BenchmarkComparison_Random_SimpleLRU(b *testing.B) {
	l, err := NewSimpleLRU(8192)
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
	b.Logf("SimpleLRU Random - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkComparison_Random_YoungOldLRU(b *testing.B) {
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
	b.Logf("YoungOldLRU Random - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

// 性能比较测试 - 频率访问模式
func BenchmarkComparison_Frequency_SimpleLRU(b *testing.B) {
	l, err := NewSimpleLRU(8192)
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
	b.Logf("SimpleLRU Frequency - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkComparison_Frequency_YoungOldLRU(b *testing.B) {
	l, err := NewYoungOldLRU(8192, 2048, 10*time.Millisecond)
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
	b.Logf("YoungOldLRU Frequency - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

// 性能比较测试 - 热点数据访问模式
func BenchmarkComparison_Hotspot_SimpleLRU(b *testing.B) {
	l, err := NewSimpleLRU(1024)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	// 80%的访问集中在20%的数据上（热点数据）
	hotKeys := make([]int64, 200)
	coldKeys := make([]int64, 800)
	for i := 0; i < 200; i++ {
		hotKeys[i] = int64(i)
	}
	for i := 0; i < 800; i++ {
		coldKeys[i] = int64(i + 200)
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < b.N; i++ {
		var key int64
		if i%10 < 8 { // 80%访问热点数据
			key = hotKeys[getRand(b)%200]
		} else { // 20%访问冷数据
			key = coldKeys[getRand(b)%800]
		}

		if i%2 == 0 {
			l.Add(key, key)
		} else {
			_, ok := l.Get(key)
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("SimpleLRU Hotspot - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkComparison_Hotspot_YoungOldLRU(b *testing.B) {
	l, err := NewYoungOldLRU(1024, 256, 50*time.Millisecond)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	// 80%的访问集中在20%的数据上（热点数据）
	hotKeys := make([]int64, 200)
	coldKeys := make([]int64, 800)
	for i := 0; i < 200; i++ {
		hotKeys[i] = int64(i)
	}
	for i := 0; i < 800; i++ {
		coldKeys[i] = int64(i + 200)
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < b.N; i++ {
		var key int64
		if i%10 < 8 { // 80%访问热点数据
			key = hotKeys[getRand(b)%200]
		} else { // 20%访问冷数据
			key = coldKeys[getRand(b)%800]
		}

		if i%2 == 0 {
			l.Add(key, key)
		} else {
			_, ok := l.Get(key)
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("YoungOldLRU Hotspot - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

// 性能比较测试 - 顺序访问模式
func BenchmarkComparison_Sequential_SimpleLRU(b *testing.B) {
	l, err := NewSimpleLRU(1024)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < b.N; i++ {
		key := int64(i % 2048) // 顺序访问，范围大于缓存大小
		if i%2 == 0 {
			l.Add(key, key)
		} else {
			_, ok := l.Get(key)
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("SimpleLRU Sequential - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkComparison_Sequential_YoungOldLRU(b *testing.B) {
	l, err := NewYoungOldLRU(1024, 256, 50*time.Millisecond)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < b.N; i++ {
		key := int64(i % 2048) // 顺序访问，范围大于缓存大小
		if i%2 == 0 {
			l.Add(key, key)
		} else {
			_, ok := l.Get(key)
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("YoungOldLRU Sequential - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

// 性能比较测试 - 纯写入操作
func BenchmarkComparison_WriteOnly_SimpleLRU(b *testing.B) {
	l, err := NewSimpleLRU(8192)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Add(int64(i), int64(i))
	}
}

func BenchmarkComparison_WriteOnly_YoungOldLRU(b *testing.B) {
	l, err := NewYoungOldLRU(8192, 2048, 100*time.Millisecond)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Add(int64(i), int64(i))
	}
}

// 性能比较测试 - 纯读取操作
func BenchmarkComparison_ReadOnly_SimpleLRU(b *testing.B) {
	l, err := NewSimpleLRU(8192)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	// 预填充数据
	for i := 0; i < 8192; i++ {
		l.Add(int64(i), int64(i))
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < b.N; i++ {
		key := getRand(b) % 8192
		_, ok := l.Get(key)
		if ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("SimpleLRU ReadOnly - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkComparison_ReadOnly_YoungOldLRU(b *testing.B) {
	l, err := NewYoungOldLRU(8192, 2048, 100*time.Millisecond)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	// 预填充数据
	for i := 0; i < 8192; i++ {
		l.Add(int64(i), int64(i))
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < b.N; i++ {
		key := getRand(b) % 8192
		_, ok := l.Get(key)
		if ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("YoungOldLRU ReadOnly - hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

// 运行性能比较测试的辅助函数
func TestRunPerformanceComparison(t *testing.T) {
	fmt.Println("\n=== LRU Performance Comparison ===")
	fmt.Println("运行以下命令进行性能测试比较：")
	fmt.Println("go test -bench=BenchmarkComparison -benchmem -count=3")
	fmt.Println("\n测试场景说明：")
	fmt.Println("1. Random: 随机访问模式 - 测试缓存在随机访问下的性能")
	fmt.Println("2. Frequency: 频率访问模式 - 测试不同频率数据的缓存效果")
	fmt.Println("3. Hotspot: 热点数据访问 - 80/20原则，测试热点数据缓存效果")
	fmt.Println("4. Sequential: 顺序访问模式 - 测试顺序访问的缓存效果")
	fmt.Println("5. WriteOnly: 纯写入操作 - 测试写入性能")
	fmt.Println("6. ReadOnly: 纯读取操作 - 测试读取性能")
	fmt.Println("\n预期结果：")
	fmt.Println("- SimpleLRU: 在简单场景下性能更好，内存开销更小")
	fmt.Println("- YoungOldLRU: 在热点数据访问场景下命中率更高，但性能开销稍大")
}
