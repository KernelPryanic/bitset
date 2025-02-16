package bitset

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		elems  []int
		expect string
	}{
		{"empty no args", []int{}, "{}"},
		{"all negatives", []int{-1, -2, -10}, "{}"},
		{"single elem", []int{1}, "{1}"},
		{"duplicates", []int{1, 1}, "{1}"},
		{"64", []int{64}, "{64}"},
		{"65", []int{65}, "{65}"},
		{"several elems", []int{1, 2, 3}, "{1..3}"},
		{"big elems", []int{100, 200, 300}, "{100 200 300}"},
		{"mixed sign", []int{1, -2, 2, -5, 3}, "{1..3}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := New(tt.elems...)
			require.Equal(t, tt.expect, bs.String())
		})
	}
}

func TestBitSet_Contains(t *testing.T) {
	bsEmpty := New()
	bsSet := New(0, 1, 2, 65, 100)

	tests := []struct {
		name     string
		bs       BitSet
		elem     int
		expected bool
	}{
		{"empty neg", bsEmpty, -1, false},
		{"empty zero", bsEmpty, 0, false},
		{"empty pos", bsEmpty, 10, false},
		{"non empty neg", bsSet, -1, false},
		{"non empty absent", bsSet, 50, false},
		{"non empty present 0", bsSet, 0, true},
		{"non empty present 65", bsSet, 65, true},
		{"non empty present 100", bsSet, 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bs.Contains(tt.elem)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestBitSet_Equal(t *testing.T) {
	tests := []struct {
		name   string
		bs1    BitSet
		bs2    BitSet
		expect bool
	}{
		{"both empty", New(), New(), true},
		{"identical small", New(1, 2), New(1, 2), true},
		{"different small", New(1, 2), New(2, 3), false},
		{"different size", New(1, 2, 65), New(1, 2), false},
		{"identical bigger", New(1, 2, 65), New(1, 2, 65), true},
		{"both large same", New(100, 200, 300), New(100, 200, 300), true},
		{"both large diff", New(100, 200, 300), New(200, 300, 400), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bs1.Equal(tt.bs2)
			require.Equal(t, tt.expect, got)
		})
	}
}

func TestBitSet_Subset(t *testing.T) {
	tests := []struct {
		name   string
		bs1    BitSet
		bs2    BitSet
		expect bool
	}{
		{"empty subset empty", New(), New(), true},
		{"empty subset non empty", New(), New(1), true},
		{"non empty subset empty", New(1), New(), false},
		{"proper subset", New(1, 2), New(1, 2, 3), true},
		{"not subset", New(1, 4), New(1, 2, 3), false},
		{"identical", New(1, 2, 3), New(1, 2, 3), true},
		{"large subset", New(100, 200), New(100, 200, 300), true},
		{"large not subset", New(100, 200, 300), New(100, 200), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bs1.Subset(tt.bs2)
			require.Equal(t, tt.expect, got)
		})
	}
}

func TestBitSet_Max(t *testing.T) {
	t.Run("negative on empty", func(t *testing.T) {
		empty := New()
		require.Equal(t, -1, empty.Max())
	})

	tests := []struct {
		name   string
		bs     BitSet
		expect int
	}{
		{"single 0", New(0), 0},
		{"single 65", New(65), 65},
		{"several", New(1, 2, 3, 62, 63, 64, 100), 100},
		{"large", New(100, 200, 300), 300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bs.Max()
			require.Equal(t, tt.expect, got)
		})
	}
}

func TestBitSet_Size(t *testing.T) {
	tests := []struct {
		name   string
		bs     BitSet
		expect int
	}{
		{"empty", New(), 0},
		{"negatives", New(-1, -5), 0},
		{"single 1", New(1), 1},
		{"single 64", New(64), 1},
		{"single 65", New(65), 1},
		{"several", New(1, 2, 3), 3},
		{"large", New(100, 200, 300), 3},
		{"range 0 to 64", func() BitSet {
			b := New()
			b.AddRange(0, 64)
			return b
		}(), 64},
		{"range 1 to 64", func() BitSet {
			b := New()
			b.AddRange(1, 64)
			return b
		}(), 63},
		{
			name: "range 0 to 576",
			bs: func() BitSet {
				b := New()
				b.AddRange(0, 576)
				return b
			}(),
			expect: 576,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bs.Size()
			require.Equal(t, tt.expect, got)
		})
	}
}

func TestBitSet_Empty(t *testing.T) {
	tests := []struct {
		name   string
		bs     BitSet
		expect bool
	}{
		{"empty", New(), true},
		{"negative only", New(-10, -5), true},
		{"range neg to 0", func() BitSet {
			b := New()
			b.AddRange(-10, 0)
			return b
		}(), true},
		{"non empty 1", New(1), false},
		{"non empty 65", New(65), false},
		{"several", New(1, 2, 3), false},
		{"large", New(100, 200, 300), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bs.Empty()
			require.Equal(t, tt.expect, got)
		})
	}
}

func TestBitSet_NextPrev(t *testing.T) {
	bs := New(0, 2, 63, 64, 100, 300)
	tests := []struct {
		name  string
		bs    BitSet
		m     int
		nextN int
		prevN int
	}{
		{"empty", New(), 1, -1, -1},
		{"empty zero", New(), 0, -1, -1},
		{"empty neg", New(), -1, -1, -1},

		{"set neg", bs, -1, 0, -1},
		{"set before 0", bs, 0, 2, -1},
		{"set between 0 and 2", bs, 1, 2, 0},
		{"set on 2", bs, 2, 63, 0},
		{"set between 2 and 63", bs, 50, 63, 2},
		{"set on 63", bs, 63, 64, 2},
		{"set on 64", bs, 64, 100, 63},
		{"set between 64 and 100", bs, 70, 100, 64},
		{"set on 100", bs, 100, 300, 64},
		{"set between 100 and 300", bs, 200, 300, 100},
		{"set on 300", bs, 300, -1, 100},
		{"past 300", bs, 400, -1, 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := tt.bs.Next(tt.m)
			p := tt.bs.Prev(tt.m)
			require.Equal(t, tt.nextN, n)
			require.Equal(t, tt.prevN, p)
		})
	}
}

func TestBitSet_Visit(t *testing.T) {
	tests := []struct {
		name   string
		bs     BitSet
		expect []int
	}{
		{"empty", New(), []int{}},
		{"single", New(0), []int{0}},
		{"several", New(1, 2, 3, 62, 63, 64), []int{1, 2, 3, 62, 63, 64}},
		{"large", New(1, 22, 333, 4444), []int{1, 22, 333, 4444}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visited := make([]int, 0)
			tt.bs.Visit(func(n int) bool {
				visited = append(visited, n)
				return false
			})
			require.Equal(t, tt.expect, visited)
		})
	}

	t.Run("abort early", func(t *testing.T) {
		bs := New(1, 2)
		count := 0
		aborted := bs.Visit(func(n int) bool {
			count++
			return n == 1
		})
		require.True(t, aborted)
		require.Equal(t, 1, count)
	})
}

func TestBitSet_VisitAll(t *testing.T) {
	bs := New(0, 2, 63, 64, 100, 300)
	visited := make([]int, 0)
	bs.VisitAll(func(n int) {
		visited = append(visited, n)
	})
	require.Equal(t, []int{0, 2, 63, 64, 100, 300}, visited)
}

func TestBitSet_Add(t *testing.T) {
	tests := []struct {
		name   string
		start  BitSet
		add    []int
		expect string
	}{
		{"add neg", New(), []int{-1}, "{}"},
		{"add single", New(), []int{1}, "{1}"},
		{"add duplicate", New(1), []int{1}, "{1}"},
		{"add new", New(1), []int{2}, "{1 2}"},
		{"add 64", New(), []int{64}, "{64}"},
		{"add large", New(), []int{100, 200, 300}, "{100 200 300}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.add {
				tt.start.Add(v)
			}
			require.Equal(t, tt.expect, tt.start.String())
		})
	}
}

func TestBitSet_Delete(t *testing.T) {
	tests := []struct {
		name   string
		start  BitSet
		del    []int
		expect string
	}{
		{"del present", New(1), []int{1}, "{}"},
		{"del neg", New(1), []int{-1}, "{1}"},
		{"del absent", New(1), []int{2}, "{1}"},
		{"del 64 from 65", New(65), []int{64}, "{65}"},
		{"del 200", New(100, 200, 300), []int{200}, "{100 300}"},
		{"del 300", New(100, 200, 300), []int{300}, "{100 200}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.del {
				tt.start.Delete(v)
			}
			require.Equal(t, tt.expect, tt.start.String())
		})
	}
}

func TestBitSet_Reset(t *testing.T) {
	tests := []struct {
		name   string
		ops    func(bs *BitSet)
		expect string
	}{
		{
			"reset empty",
			func(bs *BitSet) {},
			"{}",
		},
		{
			"reset after adding",
			func(bs *BitSet) {
				bs.Add(1)
				bs.Reset()
			},
			"{}",
		},
		{
			"reset add reset",
			func(bs *BitSet) {
				bs.Add(1)
				bs.Add(2)
				bs.Reset()
			},
			"{}",
		},
		{
			"reset add reset add",
			func(bs *BitSet) {
				bs.Add(1)
				bs.Add(2)
				bs.Reset()
				bs.Add(3)
				bs.Reset()
				bs.Add(4)
				bs.Add(2)
			},
			"{2 4}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := New()
			tt.ops(&bs)
			if tt.expect == "{}" {
				require.Equal(t, "{}", bs.String())
				return
			}
			require.Equal(t, tt.expect, bs.String())
		})
	}
}

func TestBitSet_AddRange(t *testing.T) {
	tests := []struct {
		name   string
		m, n   int
		before []int
		after  string
	}{
		{"empty range", 0, 0, nil, "{}"},
		{"empty range neg", 2, 1, nil, "{}"},
		{"neg range", -2, -1, nil, "{}"},
		{"part neg", -1, 0, nil, "{}"},
		{"simple range", 1, 10, nil, "{1..9}"},
		{"extend 64", 64, 66, nil, "{64 65}"},
		{"extend large", 1, 1000, nil, "{1..999}"},
		{"overlap existing", 1, 5, []int{1, 2, 6}, "{1..4 6}"},
		{"add on top", 50, 101, []int{100, 200}, "{50..100 200}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := New(tt.before...)
			bs.AddRange(tt.m, tt.n)
			require.Equal(t, tt.after, bs.String())
		})
	}
}

func TestBitSet_DeleteRange(t *testing.T) {
	tests := []struct {
		name   string
		m, n   int
		before []int
		after  string
	}{
		{"empty range", 0, 0, []int{1, 2, 3}, "{1..3}"},
		{"empty range neg", 2, 1, []int{1, 2, 3}, "{1..3}"},
		{"neg range", -2, -1, []int{1, 2, 3}, "{1..3}"},
		{"part neg", -1, 0, []int{0, 1}, "{0 1}"},
		{"remove part", 1, 3, []int{0, 1, 2, 3, 4}, "{0 3 4}"},
		{"remove 64", 64, 65, []int{64, 65}, "{65}"},
		{"remove big", 50, 300, []int{49, 50, 100, 200, 299, 300, 400}, "{49 300 400}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := New(tt.before...)
			bs.DeleteRange(tt.m, tt.n)
			require.Equal(t, tt.after, bs.String())
		})
	}
}

func TestBitSet_Set(t *testing.T) {
	tests := []struct {
		name string
		src  BitSet
		dst  BitSet
	}{
		{"both empty", New(), New()},
		{"dst empty src small", New(1), New()},
		{"dst empty src large", New(100, 200, 300), New()},
		{"dst small src empty", New(), New(1, 2)},
		{"dst large src empty", New(), New(64, 65, 100)},
		{"dst small src small", New(1, 2), New(2, 3)},
		{"dst large src large", New(100, 200, 300), New(50, 300, 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dstCopy := tt.dst.Copy()
			dstCopy.Set(tt.src)
			require.True(t, dstCopy.Equal(tt.src))
		})
	}
}

func TestBitSet_Copy(t *testing.T) {
	src := New(1, 2, 100, 200)
	cp := src.Copy()
	require.True(t, cp.Equal(src))

	// mutate src, ensure cp not changed
	src.Add(300)
	require.False(t, cp.Equal(src))
}

func TestAnd(t *testing.T) {
	tests := []struct {
		name   string
		a, b   BitSet
		expect string
	}{
		{"both empty", New(), New(), "{}"},
		{"a empty", New(), New(1), "{}"},
		{"b empty", New(1), New(), "{}"},
		{"overlap", New(1), New(1), "{1}"},
		{"no overlap", New(1), New(2), "{}"},
		{"partial overlap", New(1, 2), New(2, 3), "{2}"},
		{"large", New(100, 200), New(200, 300), "{200}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := And(tt.a, tt.b)
			require.Equal(t, tt.expect, res.String())
		})
	}
}

func TestBitSet_And(t *testing.T) {
	tests := []struct {
		name   string
		a, b   BitSet
		expect string
	}{
		{"both empty", New(), New(), "{}"},
		{"a empty", New(), New(1), "{}"},
		{"b empty", New(1), New(), "{}"},
		{"overlap", New(1), New(1), "{1}"},
		{"no overlap", New(1), New(2), "{}"},
		{"partial overlap", New(1, 2), New(2, 3), "{2}"},
		{"large", New(100, 200), New(200, 300), "{200}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &tt.a
			res.And(tt.b)
			require.Equal(t, tt.expect, res.String())
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		name   string
		a, b   BitSet
		expect string
	}{
		{"both empty", New(), New(), "{}"},
		{"a empty", New(), New(1), "{1}"},
		{"b empty", New(1), New(), "{1}"},
		{"same", New(1), New(1), "{1}"},
		{"no overlap", New(1), New(2), "{1 2}"},
		{"partial overlap", New(1, 2), New(2, 3), "{1..3}"},
		{"large", New(100, 200), New(200, 300), "{100 200 300}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Or(tt.a, tt.b)
			require.Equal(t, tt.expect, res.String())
		})
	}
}

func TestXor(t *testing.T) {
	tests := []struct {
		name   string
		a, b   BitSet
		expect string
	}{
		{"both empty", New(), New(), "{}"},
		{"a empty", New(), New(1), "{1}"},
		{"b empty", New(1), New(), "{1}"},
		{"equal", New(1), New(1), "{}"},
		{"2 elems equal", New(1, 2), New(1, 2), "{}"},
		{"no overlap", New(1), New(2), "{1 2}"},
		{"partial overlap", New(1, 2), New(2, 3), "{1 3}"},
		{"partial overlap trailing zero", New(1, 2, 0), New(2, 3, 0), "{1 3}"},
		{"partial overlap hundrets", New(100, 200), New(200, 300), "{100 300}"},
		{
			"20 elems no overlap",
			New(1, 100, 200, 300, 400, 500, 600, 700, 800, 900),
			New(2, 101, 202, 303, 404, 505, 606, 707, 808, 909),
			"{1 2 100 101 200 202 300 303 400 404 500 505 600 606 " +
				"700 707 800 808 900 909}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Xor(tt.a, tt.b)
			require.Equal(t, tt.expect, res.String())

			cp := tt.a.Copy()
			cp.Xor(tt.b)
			require.Equal(t, tt.expect, cp.String())
		})
	}
}

func TestAndNot(t *testing.T) {
	tests := []struct {
		name   string
		a, b   BitSet
		expect string
	}{
		{"both empty", New(), New(), "{}"},
		{"a empty", New(), New(1), "{}"},
		{"b empty", New(1), New(), "{1}"},
		{"same", New(1), New(1), "{}"},
		{"no overlap", New(1), New(2), "{1}"},
		{"partial overlap", New(1, 2), New(2, 3), "{1}"},
		{"large", New(100, 200), New(200, 300), "{100}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := AndNot(tt.a, tt.b)
			require.Equal(t, tt.expect, res.String())
		})
	}
}

func TestBitset_AndNot(t *testing.T) {
	tests := []struct {
		name   string
		a, b   BitSet
		expect string
	}{
		{"both empty", New(), New(), "{}"},
		{"a empty", New(), New(1), "{}"},
		{"b empty", New(1), New(), "{1}"},
		{"same", New(1), New(1), "{}"},
		{"no overlap", New(1), New(2), "{1}"},
		{"partial overlap", New(1, 2), New(2, 3), "{1}"},
		{"large", New(100, 200), New(200, 300), "{100}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &tt.a
			res.AndNot(tt.b)
			require.Equal(t, tt.expect, res.String())
		})
	}
}

func TestNextPow2(t *testing.T) {
	tests := []struct {
		n, expected int
	}{
		{math.MinInt, 1},
		{-1, 1},
		{0, 1},
		{1, 2},
		{2, 4},
		{3, 4},
		{4, 8},
		{1<<19 - 1, 1 << 19},
		{1 << 19, 1 << 20},
		{math.MaxInt >> 1, (math.MaxInt >> 1) + 1},
		{(math.MaxInt >> 1) + 1, math.MaxInt},
		{math.MaxInt - 1, math.MaxInt},
		{math.MaxInt, math.MaxInt},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := nextPow2(tt.n)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestBitSet_String(t *testing.T) {
	tests := []struct {
		name   string
		bs     BitSet
		expect string
	}{
		{"empty", New(), "{}"},
		{"single neg", New(-1), "{}"},
		{"single pos", New(1), "{1}"},
		{"mixed", New(1, -1), "{1}"},
		{"small range", New(0, 1, 2, 4, 5), "{0..2 4 5}"},
		{"scattered", New(0, 2, 3, 5), "{0 2 3 5}"},
		{"combined ranges", New(0, 1, 2, 3, 5, 7, 8, 9), "{0..3 5 7..9}"},
		{"single 64", New(64), "{64}"},
		{"large", New(100, 200, 300), "{100 200 300}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.bs.String()
			require.Equal(t, tt.expect, got)
		})
	}
}
