package readvarint

import (
	"encoding/binary"
	"math/rand"
	"testing"
)

func lib(p []byte) (newp []byte, val uint32) {
	x, i := binary.Uvarint(p)
	return p[i:], uint32(x)
}

func iter(p []byte) (newp []byte, val uint32) {
	var v, shift uint32
	for {
		b := p[0]
		p = p[1:]
		v |= (uint32(b) & 0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return p, v
}

func loop(p []byte) (newp []byte, val uint32) {
  v := uint(p[0])
  i := uint(1)
  for ; i<5; i++ {
    shift := uint(7*i)
    limit := uint(0x1) << shift
    if v < limit {
      break
    }
    v &^= limit
    v |= uint(p[i]) << shift
	}
	return p[i:], uint32(v)
}


func unroll(p []byte) (newp []byte, val uint32) {
	v := uint(p[0])
	if v <= 0x7F {
		return p[1:], uint32(v)
	}

	v &= 0x7F
	v |= uint(p[1]) << 7

	if v <= 0x3FFF {
		return p[2:], uint32(v)
	}

	v &= 0x3FFF
	v |= uint(p[2]) << 14

	if v <= 0x1FFFFF {
		return p[3:], uint32(v)
	}

	v &= 0x1FFFFF
	v |= uint(p[3]) << 21

	if v <= 0xFFFFFFF {
		return p[4:], uint32(v)
	}

	v &= 0xFFFFFFF
	v |= uint(p[4]) << 28

	return p[5:], uint32(v)
}

var randBytes []byte
var randInts []uint32

func init() {
	for i := 0; i < 10000; i++ {
		x := rand.Uint32()
		randInts = append(randInts, x)
    randBytes = put(randBytes, x)
	}
}

func put( b []byte, x uint32) []byte {
  size := len(b)
    b = append(b, 0, 0, 0, 0 ,0)
    l := binary.PutUvarint(b[size:], uint64(x))
    return b[:size+l]
}

func TestReadVarint(t *testing.T) {
	var data = []struct {
		name string
		readvarint    func([]byte) ([]byte, uint32)
	}{
		{"lib", lib},
		{"iter", iter},
    {"loop", loop},
		{"unroll", unroll},
	}

	ints := []uint32{
		0, 1, 2, 3, 4, 5,
		0xF, 0x10,
		0x7E, 0x7F, 0x80, 0x81,
		0xFF, 0x100,
		0x3FFE, 0x3FFF, 0x4000, 0x4001,
		0xFFFF, 0x10000,
		0x1FFFFE, 0x1FFFFF, 0x200000, 0x200001,
		0xFFFFFF, 0x1000000,
		0xFFFFFFE, 0xFFFFFFF, 0x10000000, 0x10000001,
		0xFFFFFFFE, 0xFFFFFFFF,
	}

  var bytes []byte
  for _, x := range ints {
    bytes = put(bytes, x)
  }

	for _, tt := range data {
		b := bytes
		for _, want := range ints {
			var got uint32
			b, got = tt.readvarint(b)
			if got != want {
				t.Errorf("%s: %d(%x) != %d(%x)", tt.name, got, got, want, want)
			}
		}
	}
}

func BenchmarkLib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes := randBytes
		for _, want := range randInts {
			var got uint32
			bytes, got = lib(bytes)
			if got != want {
				b.Fatalf("%d(%x) != %d(%x)", got, got, want, want)
			}
		}
	}
}

func BenchmarkIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes := randBytes
		for _, want := range randInts {
			var got uint32
			bytes, got = iter(bytes)
			if got != want {
				b.Fatalf("%d(%x) != %d(%x)", got, got, want, want)
			}
		}
	}
}

func BenchmarkLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes := randBytes
		for _, want := range randInts {
			var got uint32
			bytes, got = loop(bytes)
			if got != want {
				b.Fatalf("%d(%x) != %d(%x)", got, got, want, want)
			}
		}
	}
}

func BenchmarkUnroll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes := randBytes
		for _, want := range randInts {
			var got uint32
			bytes, got = unroll(bytes)
			if got != want {
				b.Fatalf("%d(%x) != %d(%x)", got, got, want, want)
			}
		}
	}
}
