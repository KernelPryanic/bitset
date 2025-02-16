package bitset

import (
	"fmt"
	"math"
	"math/bits"
	"strings"
)

const (
	bpw      = 64         // bits per word
	maxw     = 1<<bpw - 1 // maximum value of a word
	shift    = 6          // 1<<6 == 64, to be used as multiplier/divisor by 64
	div64rem = 63         // remainder of division by 64 when n&div64rem is used
)

// Word is a convenience alias.
type Word = uint64

// BitSet is a set of non-negative integers represented as a slice of uint64 words,
// where each bit i in word w corresponds to the integer 64*n + i.
// The words are kept in ascending order, and the set is trimmed
// to remove trailing zero words.
type BitSet []uint64

// New creates a new set with the given non-negative elements.
// If all n are negative, an empty set is created. The elements are stored in
// ascending order. The zero value of BitSet is an empty set.
func New(n ...int) BitSet {
	if len(n) == 0 {
		return BitSet{}
	}
	maxElem := n[0]
	for _, e := range n {
		if e > maxElem {
			maxElem = e
		}
	}
	if maxElem < 0 {
		return BitSet{}
	}
	s := make(BitSet, (maxElem>>shift)+1)
	for _, e := range n {
		if e >= 0 {
			s[e>>shift] |= 1 << uint(e&div64rem)
		}
	}
	return s
}

// Reset resets the set without reallocation.
func (bs *BitSet) Reset() {
	for i := range *bs {
		(*bs)[i] = 0
	}
	*bs = (*bs)[:0]
}

// Contains tells if n is in the set.
func (bs BitSet) Contains(n int) bool {
	if n < 0 {
		return false
	}
	i := n >> shift
	if i >= len(bs) {
		return false
	}
	return (bs[i] & (1 << uint(n&div64rem))) != 0
}

// Equal tells if bs and other are equal.
func (bs BitSet) Equal(other BitSet) bool {
	if len(bs) != len(other) {
		return false
	}
	for i := 0; i < len(bs); i++ {
		if bs[i] != other[i] {
			return false
		}
	}
	return true
}

// Subset tells if bs is a subset of other.
func (bs BitSet) Subset(other BitSet) bool {
	bsLen := len(bs)
	if bsLen > len(other) {
		return false
	}
	for i := 0; i < bsLen; i++ {
		if bs[i]&^other[i] != 0 {
			return false
		}
	}
	return true
}

// Max returns the maximum element of the bitset.
// If the set is empty, -1 is returned.
func (bs BitSet) Max() int {
	bsLen := len(bs)
	if bsLen == 0 {
		return -1
	}
	i := bsLen - 1
	return (i << shift) + bits.Len64(bs[i]) - 1
}

// Size returns the number of elements in the set.
func (bs BitSet) Size() int {
	size := 0
	for _, word := range bs {
		size += bits.OnesCount64(word)
	}
	return size
}

// Empty tells if the set is empty.
func (bs BitSet) Empty() bool {
	return len(bs) == 0
}

// Next returns the next element n, n > m, in the set,
// or -1 if there is no such element.
func (bs BitSet) Next(m int) int {
	l := len(bs)
	if l == 0 {
		return -1
	}
	if m < 0 {
		if (bs[0] & 1) != 0 {
			return 0
		}
		m = 0
	}
	i := m >> shift
	if i >= l {
		return -1
	}
	t := 1 + uint(m&div64rem)
	w := bs[i] >> t << t // zero out bits for numbers ≤ m
	for i < l-1 && w == 0 {
		i++
		w = bs[i]
	}
	if w == 0 {
		return -1
	}
	return (i << shift) + bits.TrailingZeros64(w)
}

// Prev returns the previous element n, n < m, in the set,
// or -1 if there is no such element.
func (bs BitSet) Prev(m int) int {
	l := len(bs)
	if l == 0 || m <= 0 {
		return -1
	}
	lastIdx := l - 1
	maxPossible := (lastIdx << shift) + bits.Len64(bs[lastIdx]) - 1
	if m > maxPossible {
		return maxPossible
	}
	i := m >> shift
	t := bpw - uint(m&div64rem)
	if i >= l {
		i = l - 1
		t = bpw
	}
	w := bs[i] << t >> t // zero out bits >= m
	for i > 0 && w == 0 {
		i--
		w = bs[i]
	}
	if w == 0 {
		return -1
	}
	return (i << shift) + bits.Len64(w) - 1
}

// Visit calls the do function for each element of s in numerical order.
// If do returns true, Visit returns immediately, skipping any remaining
// elements, and returns true. It is safe for do to add or delete
// elements e, e ≤ n. The behavior of Visit is undefined if do changes
// the set in any other way.
func (bs BitSet) Visit(do func(n int) bool) (aborted bool) {
	for i, l := 0, len(bs); i < l; i++ {
		w := bs[i]
		if w == 0 {
			continue
		}
		n := i << shift
		for w != 0 {
			b := bits.TrailingZeros64(w)
			n += b
			if do(n) {
				return true
			}
			n++
			w >>= (b + 1)
			for w&1 != 0 {
				if do(n) {
					return true
				}
				n++
				w >>= 1
			}
		}
	}
	return false
}

// VisitAll calls do function for each element of s in numerical order.
func (bs BitSet) VisitAll(do func(n int)) {
	bs.Visit(func(n int) bool {
		do(n)
		return false
	})
}

// bitMask returns a bit div64rem with bits set from m to n inclusive, 0 ≤ m ≤ n < bpw.
func bitMask(m, n int) uint64 {
	return maxw >> uint(bpw-1-(n-m)) << uint(m)
}

// nextPow2 returns the smallest power of two p such that p > n, or math.MaxInt if overflow.
func nextPow2(n int) int {
	if n <= 0 {
		return 1
	}
	k := bits.Len64(uint64(n))
	if k < bits.UintSize-1 {
		return 1 << uint(k)
	}
	return math.MaxInt
}

// newCap suggests a new increased capacity when growing a slice to length n.
func newCap(n, prevCap int) int {
	return max(n, nextPow2(prevCap))
}

// resize changes the capacity of *bs to hold at least n elements.
// If n is less than the current length of *bs, the set is truncated.
func (bs *BitSet) resize(n int) {
	if cap(*bs) < n {
		newData := make(BitSet, n, newCap(n, cap(*bs)))
		copy(newData, *bs)
		*bs = newData
	}
	for i := len(*bs) - 1; i >= n && i >= 0; i-- {
		(*bs)[i] = 0
	}
	*bs = (*bs)[:n]
}

// trim slices *bs by removing all trailing words equal to zero.
func (bs *BitSet) trim() {
	i := len(*bs) - 1
	for i >= 0 && (*bs)[i] == 0 {
		i--
	}
	*bs = (*bs)[:i+1]
}

// Set replaces the contents of *bs with other.
func (bs *BitSet) Set(other BitSet) {
	*bs = make(BitSet, len(other))
	copy(*bs, other)
}

// Copy creates a new set that is a copy of bs.
func (bs BitSet) Copy() BitSet {
	if len(bs) == 0 {
		return BitSet{}
	}
	s := make(BitSet, len(bs))
	copy(s, bs)
	return s
}

// Add adds n to bs (no-op if n < 0).
func (bs *BitSet) Add(n int) {
	if n < 0 {
		return
	}
	i := n >> shift
	if i >= len(*bs) {
		bs.resize(i + 1)
	}
	(*bs)[i] |= 1 << uint(n&div64rem)
}

// Delete removes n from bs (no-op if n < 0 or not present).
func (bs *BitSet) Delete(n int) {
	if n < 0 {
		return
	}
	i := n >> shift
	if i >= len(*bs) {
		return
	}
	(*bs)[i] &^= 1 << uint(n&div64rem)
	bs.trim()
}

// AddRange adds all integers from m to n-1 to bs (no-op if m>=n).
func (bs *BitSet) AddRange(m, n int) {
	if n < 1 || m >= n {
		return
	}
	m = max(0, m)
	n-- // now we add from m..n inclusive
	low, high := m>>shift, n>>shift
	if high >= len(*bs) {
		bs.resize(high + 1)
	}
	if low == high {
		(*bs)[low] |= bitMask(m&div64rem, n&div64rem)
		return
	}
	(*bs)[low] |= bitMask(m&div64rem, bpw-1)
	for i := low + 1; i < high; i++ {
		(*bs)[i] = maxw
	}
	(*bs)[high] |= bitMask(0, n&div64rem)
}

// DeleteRange removes all integers from m to n-1 (no-op if m>=n).
func (bs *BitSet) DeleteRange(m, n int) {
	if n < 1 || m >= n {
		return
	}
	m = max(0, m)
	n-- // remove m..n inclusive
	low, high := m>>shift, n>>shift
	if low >= len(*bs) {
		return
	}
	if high >= len(*bs) {
		high = len(*bs) - 1
		n = bpw - 1
	}
	if low == high {
		(*bs)[low] &^= bitMask(m&div64rem, n&div64rem)
		bs.trim()
		return
	}
	(*bs)[low] &^= bitMask(m&div64rem, bpw-1)
	for i := low + 1; i < high; i++ {
		(*bs)[i] = 0
	}
	(*bs)[high] &^= bitMask(0, n&div64rem)
	bs.trim()
}

// And creates a new set that consists of all elements in both s1 and s2.
func And(s1, s2 BitSet) BitSet {
	bsLen, otherLen := len(s1), len(s2)
	if bsLen == 0 || otherLen == 0 {
		return BitSet{}
	}
	minLen := min(bsLen, otherLen) - 1
	for minLen >= 0 && s1[minLen]&s2[minLen] == 0 {
		minLen--
	}
	s := make(BitSet, minLen+1)
	for i := 0; i <= minLen; i++ {
		s[i] = s1[i] & s2[i]
	}
	return s
}

// And keeps only bits set in both *bs and other.
func (bs *BitSet) And(other BitSet) {
	minLen := min(len(*bs), len(other))
	for i := 0; i < minLen; i++ {
		(*bs)[i] &= other[i]
	}
	for i := minLen; i < len(*bs); i++ {
		(*bs)[i] = 0
	}
	bs.trim()
}

// Or creates a new set that contains all elements in s1 or s2.
func Or(s1, s2 BitSet) BitSet {
	if len(s1) < len(s2) {
		s1, s2 = s2, s1 // swap to make s1 the longer set
	}
	bsLen, otherLen := len(s1), len(s2)
	if bsLen == 0 {
		return BitSet{}
	}
	if otherLen == 0 { // s2 is empty, return a copy of s1 with trailing zeros removed
		last := bsLen - 1
		for last >= 0 && s1[last] == 0 {
			last--
		}
		if last < 0 {
			return BitSet{}
		}
		s := make(BitSet, last+1)
		copy(s, s1[:last+1])
		return s
	}
	n := bsLen - 1
	for n >= 0 { // find the last non-zero word in s1 or s2
		if s1[n] != 0 || n < otherLen && s2[n] != 0 {
			break
		}
		n--
	}
	if n < 0 {
		return BitSet{}
	}
	s := make(BitSet, n+1)
	for i := 0; i <= n; i++ {
		if i < otherLen {
			s[i] = s1[i] | s2[i]
		} else {
			s[i] = s1[i]
		}
	}
	return s
}

// Or sets bits that are set in either *bs or other.
func (bs *BitSet) Or(other BitSet) {
	if len(other) > len(*bs) {
		bs.resize(len(other))
	}
	for i := 0; i < len(other); i++ {
		(*bs)[i] |= other[i]
	}
	bs.trim()
}

// Xor creates a new set that contains all elements in s1 or s2 but not both.
func Xor(s1, s2 BitSet) BitSet {
	if len(s1) < len(s2) {
		s1, s2 = s2, s1 // swap to make s1 the longer set
	}
	bsLen, otherLen := len(s1), len(s2)
	n := bsLen - 1
	for n >= 0 { // find the last non-zero word in s1 or s2
		if s1[n] != 0 || n < otherLen && s2[n] != 0 {
			break
		}
		n--
	}
	if n < 0 {
		return BitSet{}
	}
	s := make(BitSet, n+1)
	for i := 0; i <= n; i++ {
		if i < otherLen {
			s[i] = s1[i] ^ s2[i]
		} else {
			s[i] = s1[i]
		}
	}
	return s
}

// Xor toggles bits that are set in either *bs or other but not both.
func (bs *BitSet) Xor(other BitSet) {
	if len(other) > len(*bs) {
		bs.resize(len(other))
	}
	if len(other) < 8 {
		for i := range other {
			(*bs)[i] ^= other[i]
		}
		bs.trim()
		return
	}

	b := *bs
	for ; len(other) > 7; b, other = b[8:], other[8:] {
		(b)[0] ^= other[0]
		(b)[1] ^= other[1]
		(b)[2] ^= other[2]
		(b)[3] ^= other[3]
		(b)[4] ^= other[4]
		(b)[5] ^= other[5]
		(b)[6] ^= other[6]
		(b)[7] ^= other[7]
	}
	for i := range other {
		b[i] ^= other[i]
	}
	bs.trim()
}

// AndNot creates a new set that consists of all elements in s1 but not in s2.
func AndNot(s1, s2 BitSet) BitSet {
	bsLen, otherLen := len(s1), len(s2)
	if bsLen == 0 {
		return BitSet{}
	}
	n := bsLen - 1
	for n >= 0 && s1[n] == 0 { // find the last non-zero word in s1, we care only about s1 because AND NOT is always zero when s1 is zero
		n--
	}
	if n < 0 {
		return BitSet{}
	}
	s := make(BitSet, n+1)
	for i := 0; i <= n; i++ {
		if i < otherLen {
			s[i] = s1[i] &^ s2[i]
		} else {
			s[i] = s1[i]
		}
	}
	return s
}

// AndNot removes bits that are set in other from *bs.
func (bs *BitSet) AndNot(other BitSet) {
	minLen := min(len(*bs), len(other))
	for i := 0; i < minLen; i++ {
		(*bs)[i] &^= other[i]
	}
	bs.trim()
}

// writeRange appends either "", "a", "a b" or "a..b" to buf.
func writeRange(buf *strings.Builder, a, b int) {
	switch {
	case a > b:
		return
	case a == b:
		fmt.Fprintf(buf, "%d", a)
	case a+1 == b:
		fmt.Fprintf(buf, "%d %d", a, b)
	default:
		fmt.Fprintf(buf, "%d..%d", a, b)
	}
}

// String returns a string representation of the set.
// The elements are listed in ascending order.
//
// Example: {0 2 4..7 9 11 13 15}
func (bs BitSet) String() string {
	buf := new(strings.Builder)
	buf.WriteByte('{')
	a, b := -1, -2
	first := true
	bs.Visit(func(n int) bool {
		if n == b+1 {
			b++
			return false
		}
		if first && a <= b {
			first = false
		} else if a <= b {
			buf.WriteByte(' ')
		}
		writeRange(buf, a, b)
		a, b = n, n
		return false
	})
	if !first && a <= b {
		buf.WriteByte(' ')
	}
	writeRange(buf, a, b)
	buf.WriteByte('}')
	return buf.String()
}
