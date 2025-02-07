// bench_test.go
package bitset

import (
	"strconv"
	"testing"
)

// setupBenchmarkSets creates BitSets of various sizes for benchmarking
func setupBenchmarkSets() (BitSet, BitSet) {
	small := New()
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			small.Add(i)
		}
	}

	large := New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			large.Add(i)
		}
	}

	return small, large
}

func BenchmarkNew(b *testing.B) {
	elements := []int{1, 10, 100, 1000, 10000}
	for _, size := range elements {
		nums := make([]int, size)
		for i := 0; i < size; i++ {
			nums[i] = i
		}
		b.Run("size "+strconv.Itoa(size), func(b *testing.B) {
			for b.Loop() {
				New(nums...)
			}
		})
	}
}

func BenchmarkBitSet_Reset(b *testing.B) {
	small, large := setupBenchmarkSets()

	b.Run("small set", func(b *testing.B) {
		for b.Loop() {
			small.Reset()
		}
	})

	b.Run("large set", func(b *testing.B) {
		for b.Loop() {
			large.Reset()
		}
	})
}

func BenchmarkBitSet_Add(b *testing.B) {
	scenarios := []struct {
		name string
		n    int
	}{
		{"small", 10},
		{"medium", 1000},
		{"large", 100000},
	}

	for _, sc := range scenarios {
		b.Run(sc.name, func(b *testing.B) {
			bs := New()
			b.ResetTimer()
			for b.Loop() {
				bs.Add(bs.Size() % sc.n)
			}
		})
	}
}

func BenchmarkBitSet_Contains(b *testing.B) {
	small, large := setupBenchmarkSets()

	b.Run("small present", func(b *testing.B) {
		for b.Loop() {
			small.Contains(50)
		}
	})

	b.Run("small absent", func(b *testing.B) {
		for b.Loop() {
			small.Contains(51)
		}
	})

	b.Run("large present", func(b *testing.B) {
		for b.Loop() {
			large.Contains(5000)
		}
	})

	b.Run("large absent", func(b *testing.B) {
		for b.Loop() {
			large.Contains(5001)
		}
	})
}

func BenchmarkBitSet_AddRange(b *testing.B) {
	scenarios := []struct {
		name  string
		start int
		end   int
	}{
		{"small range", 0, 100},
		{"medium range", 0, 1000},
		{"large range", 0, 10000},
	}

	for _, sc := range scenarios {
		b.Run(sc.name, func(b *testing.B) {
			for b.Loop() {
				bs := New()
				bs.AddRange(sc.start, sc.end)
			}
		})
	}
}

func BenchmarkBitSet_DeleteRange(b *testing.B) {
	scenarios := []struct {
		name  string
		start int
		end   int
	}{
		{"small range", 0, 100},
		{"medium range", 0, 1000},
		{"large range", 0, 10000},
	}

	for _, sc := range scenarios {
		b.Run(sc.name, func(b *testing.B) {
			for b.Loop() {
				bs := New()
				bs.AddRange(sc.start, sc.end)
				bs.DeleteRange(sc.start, sc.end)
			}
		})
	}
}

func BenchmarkAnd(b *testing.B) {
	small1, small2 := New(1, 2, 3, 4, 5), New(3, 4, 5, 6, 7)
	large1, large2 := New(), New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			large1.Add(i)
		}
		if i%3 == 0 {
			large2.Add(i)
		}
	}

	b.Run("small sets", func(b *testing.B) {
		for b.Loop() {
			_ = And(small1, small2)
		}
	})

	b.Run("large sets", func(b *testing.B) {
		for b.Loop() {
			_ = And(large1, large2)
		}
	})
}

func BenchmarkBitSet_And(b *testing.B) {
	small1, small2 := New(1, 2, 3, 4, 5), New(3, 4, 5, 6, 7)
	large1, large2 := New(), New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			large1.Add(i)
		}
		if i%3 == 0 {
			large2.Add(i)
		}
	}

	b.Run("small sets", func(b *testing.B) {
		for b.Loop() {
			small1.And(small2)
		}
	})

	b.Run("large sets", func(b *testing.B) {
		for b.Loop() {
			large1.And(large2)
		}
	})
}

func BenchmarkOr(b *testing.B) {
	small1, small2 := New(1, 2, 3, 4, 5), New(3, 4, 5, 6, 7)
	large1, large2 := New(), New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			large1.Add(i)
		}
		if i%3 == 0 {
			large2.Add(i)
		}
	}

	b.Run("small sets", func(b *testing.B) {
		for b.Loop() {
			_ = Or(small1, small2)
		}
	})

	b.Run("large sets", func(b *testing.B) {
		for b.Loop() {
			_ = Or(large1, large2)
		}
	})
}

func BenchmarkBitSet_Or(b *testing.B) {
	small1, small2 := New(1, 2, 3, 4, 5), New(3, 4, 5, 6, 7)
	large1, large2 := New(), New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			large1.Add(i)
		}
		if i%3 == 0 {
			large2.Add(i)
		}
	}

	b.Run("small sets", func(b *testing.B) {
		for b.Loop() {
			small1.Or(small2)
		}
	})

	b.Run("large sets", func(b *testing.B) {
		for b.Loop() {
			large1.Or(large2)
		}
	})
}

func BenchmarkXor(b *testing.B) {
	small1, small2 := New(1, 2, 3, 4, 5), New(3, 4, 5, 6, 7)
	large1, large2 := New(), New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			large1.Add(i)
		}
		if i%3 == 0 {
			large2.Add(i)
		}
	}

	b.Run("small sets", func(b *testing.B) {
		for b.Loop() {
			_ = Xor(small1, small2)
		}
	})

	b.Run("large sets", func(b *testing.B) {
		for b.Loop() {
			_ = Xor(large1, large2)
		}
	})
}

func BenchmarkBitSet_Xor(b *testing.B) {
	small1, small2 := New(1, 2, 3, 4, 5), New(3, 4, 5, 6, 7)
	large1, large2 := New(), New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			large1.Add(i)
		}
		if i%3 == 0 {
			large2.Add(i)
		}
	}

	b.Run("small sets", func(b *testing.B) {
		for b.Loop() {
			small1.Xor(small2)
		}
	})

	b.Run("large sets", func(b *testing.B) {
		for b.Loop() {
			large1.Xor(large2)
		}
	})
}

// --- NEW: Benchmarks for AND NOT ---

func BenchmarkAndNot(b *testing.B) {
	small1, small2 := New(1, 2, 3, 4, 5), New(3, 4, 5, 6, 7)
	large1, large2 := New(), New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			large1.Add(i)
		}
		if i%3 == 0 {
			large2.Add(i)
		}
	}

	type benchCase struct {
		name string
		s1   BitSet
		s2   BitSet
	}

	cases := []benchCase{
		{"small sets", small1, small2},
		{"large sets", large1, large2},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for b.Loop() {
				_ = AndNot(c.s1, c.s2)
			}
		})
	}
}

func BenchmarkBitSet_AndNot(b *testing.B) {
	small1, small2 := New(1, 2, 3, 4, 5), New(3, 4, 5, 6, 7)
	large1, large2 := New(), New()
	for i := 0; i < 10000; i++ {
		if i%2 == 0 {
			large1.Add(i)
		}
		if i%3 == 0 {
			large2.Add(i)
		}
	}

	type benchCase struct {
		name string
		s1   BitSet
		s2   BitSet
	}

	cases := []benchCase{
		{"small sets", small1, small2},
		{"large sets", large1, large2},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for b.Loop() {
				c.s1.AndNot(c.s2)
			}
		})
	}
}

// -----------------------------------

func BenchmarkBitSet_Copy(b *testing.B) {
	small, large := setupBenchmarkSets()

	b.Run("small set", func(b *testing.B) {
		for b.Loop() {
			small.Copy()
		}
	})

	b.Run("large set", func(b *testing.B) {
		for b.Loop() {
			large.Copy()
		}
	})
}

func BenchmarkBitSet_String(b *testing.B) {
	small, large := setupBenchmarkSets()

	b.Run("small set", func(b *testing.B) {
		for b.Loop() {
			_ = small.String()
		}
	})

	b.Run("large set", func(b *testing.B) {
		for b.Loop() {
			_ = large.String()
		}
	})
}

func BenchmarkBitSet_Visit(b *testing.B) {
	small, large := setupBenchmarkSets()
	dummy := 0 // Used to prevent compiler optimizations

	b.Run("small set", func(b *testing.B) {
		for b.Loop() {
			small.Visit(func(n int) bool {
				dummy += n
				return false
			})
		}
	})

	b.Run("large set", func(b *testing.B) {
		for b.Loop() {
			large.Visit(func(n int) bool {
				dummy += n
				return false
			})
		}
	})
}
