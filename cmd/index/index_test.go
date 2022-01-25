package main

import (
	"bytes"
	"strings"
	"testing"
)

const position = 25
const repeat = 10000

var bdata = []byte(strings.Repeat("1", position) + "\t" + "dsahfk/sdfhkd/sj")
var sdata = strings.Repeat("1", position) + "\t" + "dsahfk/sdfhkd/sj"

func BenchmarkBytesIndexByteStdLib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < repeat; i++ {
			pos := bytes.IndexByte(bdata, '\t')
			if pos != position {
				b.Fatalf("Wrong position %d, %d expected", pos, position)
			}
		}
	}
}

//func BenchmarkBytesIndexByteStdLibV2(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		for i := 0; i < repeat; i++ {
//			pos := bytes.Index(data, []byte("\t"))
//			if pos != position {
//				b.Fatalf("Wrong position %d, %d expected", pos, position)
//			}
//		}
//	}
//}
//
//func BenchmarkCustomBytesIndexByteLoop(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		for i := 0; i < repeat; i++ {
//			pos := CustomBytesIndexByte(data, '\t')
//			if pos != position {
//				b.Fatalf("Wrong position %d, %d expected", pos, position)
//			}
//		}
//	}
//}

func BenchmarkStringIndexByteStdLib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < repeat; i++ {
			pos := strings.IndexByte(sdata, '\t')
			if pos != position {
				b.Fatalf("Wrong position %d, %d expected", pos, position)
			}
		}
	}
}

//func BenchmarkStringIndexByteStdLibV2(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		for i := 0; i < repeat; i++ {
//			pos := strings.Index(string(data), "\t")
//			if pos != position {
//				b.Fatalf("Wrong position %d, %d expected", pos, position)
//			}
//		}
//	}
//}
//
//func BenchmarkCustomStringIndexByteLoop(b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		for i := 0; i < repeat; i++ {
//			pos := CustomStringIndexByte(string(data), '\t')
//			if pos != position {
//				b.Fatalf("Wrong position %d, %d expected", pos, position)
//			}
//		}
//	}
//}

func CustomBytesIndexByte(bb []byte, b byte) int {
	pos := -1
	for j, c := range bb {
		if c == b {
			pos = j
			break
		}
	}
	return pos
}

func CustomStringIndexByte(s string, r rune) int {
	pos := -1
	for j, c := range s {
		if c == r {
			pos = j
			break
		}
	}
	return pos
}

func BenchmarkStringCountStdLib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < repeat; i++ {
			n := strings.Count(sdata, "/")
			if n != 2 {
				b.Fatalf("Wrong count %d, %d expected", n, 2)
			}
		}
	}
}

func BenchmarkStringCountCustom(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < repeat; i++ {
			n := StringCountCustom(sdata, '/')
			if n != 2 {
				b.Fatalf("Wrong count %d, %d expected", n, 2)
			}
		}
	}
}

func StringCountCustom(s string, c byte) int {
	n := 0
	for {
		i := strings.IndexByte(s, c)
		if i == -1 {
			return n
		}
		n++
		s = s[i+1:]
	}
	return n
}
