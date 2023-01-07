package tablegen

import (
	"fmt"
	"strings"
)

type byteSet struct {
	data map[byte]bool
}

func newByteSet(values... byte) byteSet {
	res := byteSet{ data: map[byte]bool{}}
	for _, v := range values {
		res.Add(v)
	}
	return res
}

func (s byteSet) Contains(value byte) bool {
	_, ok := s.data[value]
	return ok
}

func (s byteSet) ToSlice() []byte {
	res := []byte{}
	for k := range s.data {
		res = append(res, k)
	}
	return res
}

func (s byteSet) Add(value byte) {
	s.data[value] = true
	return
}

func (s byteSet) Len() int {
	return len(s.data)
}

func (s byteSet) String() string {
	sb := strings.Builder{}
	sb.WriteString("byteSet[")

	first := true
	for k := range s.data {
		if !first {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%02X", k))
		first = false
	}
	sb.WriteByte(']')
	return sb.String()
}
