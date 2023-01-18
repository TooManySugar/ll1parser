//
// Fixed size byteSet - 32 byte
//
// Add, Remove, Contains, Union, Equal have O(1) complexity
// Bytes, String, Len - O(n) where n = count of bytes in set
//     Maybe they can be faster I'm not that smart
//
// This implementation faster compared to Tree-bases sets on
// Add, Remove, Contains, Union, Equal and as fast on other tasks
//
// About space complexity
// If we assume 'perfect' Tree-based set exists it's node size would be
// 1 byte making this byteset worse on spacewise when working on sets with
// items count < 32.
// But realistically such trees can not exist.
// Smallest node size for byte set tree would be 3 bytes (basic BST):
//     1 byte  - Value itself
//     2 bytes - Indexes of L and R nodes in node array with max size of 256
// So this implementation would be preferable spacewise then working on sets
// with items count > 10.
// Nonetheless you greedy of 32 bytes of RAM on byteSet?
//
package byteset

import (
	"fmt"
	// bytes used for buffering in ByteSet.String function
	"bytes"
)

type ByteSet struct {
	data [4]uint64
}

func New(values... byte) ByteSet {
	res := ByteSet{}
	for _, v := range values {
		res.Add(v)
	}
	return res
}

func (s *ByteSet) Add(value byte) {
	i := value / 64
	mask := uint64(1) << (value % 64)
	s.data[i] = s.data[i] | mask
}

func (s *ByteSet) Remove(value byte) {
	i := value / 64
	mask := uint64(1) << (value % 64)
	mask = ^mask
	s.data[i] = s.data[i] & mask
}

func (s ByteSet) Contains(value byte) bool {
	i := int(value) / 64
	offset := int(value) % 64
	return (s.data[i] >> offset) & 1 == 1
}

// Peter Wegner's / Derrick Lehmer's method of set bits counting
// also known as Brian Kernighan's
// 1 loop cycle per set bit making it as efficient as tree traverse
func popCount(v uint64) int {
	var c int
	for c = 0; v != 0; c++ {
		v &= v - 1
	}
	return c
}

func (s ByteSet) Len() int {
	res := 0
	res += popCount(s.data[0])
	res += popCount(s.data[1])
	res += popCount(s.data[2])
	res += popCount(s.data[3])
	return res
}

var lookupTable = [64]uint8{
	 0,  1, 56,  2, 57, 16,  3, 47, 61, 58, 17, 41, 32,  4, 36, 48,
	62, 14, 59, 30, 12, 18, 42, 20, 44, 33, 27,  5, 37,  8, 22, 49,
	63, 55, 15, 46, 60, 40, 31, 35, 13, 29, 11, 19, 43, 26,  7, 21,
	54, 45, 39, 34, 28, 10, 25,  6, 53, 38,  9, 24, 52, 23, 51, 50}

// Searches Log base 2 from power of 2 numbers with O(1)
// using De Bruijn Sequence
//
// For how to get this hexadecimal nonsense and lookup table refer:
//     https://github.com/TooManySugar/debruijnbtwon/tree/master/examples
//
// Original implementation (for powers of 2 up to 5) taken from:
//     https://graphics.stanford.edu/~seander/bithacks.html#IntegerLogDeBruijn
//
func log2pow2n(pow2n uint64) (n uint8) {
	return lookupTable[(pow2n * 0x037515ED33963F09) >> 58]
}

// Modified version of popCount to retrive set bits and covert them to
// corresponding bytes for further manipulations on them
//
func realBytesF(v uint64, o byte, f func(byte)) {
	for v != 0 {
		i := v
		v &= v - 1
		// isolating right most bit of v (smallest byte in data[*])
		i = i - v
		// adding to it position offset
		b := o + log2pow2n(i)
		f(b)
	}
}

func (s ByteSet) walkBytes(f func(byte)) {
	realBytesF(s.data[0], 0x00, f)
	realBytesF(s.data[1], 0x40, f)
	realBytesF(s.data[2], 0x80, f)
	realBytesF(s.data[3], 0xC0, f)
}

// Returns slice of bytes in set in increasing order
func (s ByteSet) Bytes() []byte {
	res := []byte{}
	f := func(b byte) {
		res = append(res, b)
	}

	s.walkBytes(f)

	return res
}

// Returns a new ByteSet representing the union
func (s ByteSet) Union(other ByteSet) ByteSet {
	return ByteSet {
		[4]uint64{
			s.data[0] | other.data[0],
			s.data[1] | other.data[1],
			s.data[2] | other.data[2],
			s.data[3] | other.data[3],
		},
	}
}

// Returns true if other is the same set
func (s ByteSet) Equal(other ByteSet) bool {
	return s.data[0] == other.data[0] &&
	       s.data[1] == other.data[1] &&
	       s.data[2] == other.data[2] &&
	       s.data[3] == other.data[3]
}

// fmt.Stringer implementation
func (s ByteSet) String() string {
	sb := bytes.Buffer{}
	fmt.Fprint(&sb, "[")

	c := 0
	f := func(b byte) {
		t, _ := fmt.Fprintf(&sb, "%02X ", b)
		c += t
	}

	s.walkBytes(f)

	if c == 0 {
		return "[]"
	}

	b := sb.Bytes()
	b[len(b)-1] = ']'
	return string(b)
}
