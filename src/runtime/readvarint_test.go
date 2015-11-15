package readvarint_test

import (
	"encoding/binary"
	"testing"
  "math/rand"
)

func lib(p []byte) (newp []byte, val uint32) {
	x, i := binary.Uvarint(p)
	return p[i:], uint32(x)
}

func plain(p []byte) (newp []byte, val uint32) {
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

func clean(p []byte) (newp []byte, val uint32) {
	var v uint32

  v |= uint32(p[0]&0x7F) << (0 * 7)
  if p[0] < 128 {
    return p[0+1:], v
  }

  v |= uint32(p[1]&0x7F) << (1 * 7)
  if p[1] < 128 {
    return p[1+1:], v
  }

  v |= uint32(p[2]&0x7F) << (2 * 7)
  if p[2] < 128 {
    return p[2+1:], v
  }

  v |= uint32(p[3]&0x7F) << (3 * 7)
  if p[3] < 128 {
    return p[3+1:], v
  }

  v |= uint32(p[4]&0x7F) << (4 * 7)
  if p[4] < 128 {
    return p[4+1:], v
  }

	return p[5:], v
}

func magic(p []byte) (newp []byte, val uint32) {
	v := uint32(p[0] & 0x7F)
	if p[0] < 128 {
		return p[1:], v
	}

	v |= uint32(p[1]&0x7F) << 7
	if p[1] < 128 {
		return p[2:], v
	}

	v |= uint32(p[2]&0x7F) << 14
	if p[2] < 128 {
		return p[3:], v
	}

	v |= uint32(p[3]&0x7F) << 21
	if p[3] < 128 {
		return p[4:], v
	}

	return p[5:], v | (uint32(p[4]&0x7F) << 28)
}

var randBytes []byte
var randInts []uint32

func init() {
  var scratch [5]byte
  for i:=0; i<1000; i++ {
    x := rand.Uint32()
    randInts = append(randInts, x)
    l := binary.PutUvarint(scratch[:], uint64(x))
    randBytes = append(randBytes, scratch[:l]...)
  }
}

func ints2varint(ints []uint32) []byte {
	var b []byte
	scratch := make([]byte, 5)
	for _, i := range ints {
		l := binary.PutUvarint(scratch, uint64(i))
		varint := scratch[:l]
		b = append(b, varint...)
	}

	return append(b, 0)
}

func TestReadVarint(t *testing.T) {
	var data = []struct {
		name string
		f    func([]byte) ([]byte, uint32)
	}{
		{"lib", lib},
		{"plain", plain},
		{"clean", clean},
		{"magic", magic},
	}

	ints := []uint32{0, 1, 2, 3, 4, 5, 0xFE, 0xFF, 0x100, 0xFFFE, 0xFFFF, 0x10000, 0xFFFFFFFF}

	bytes := ints2varint(ints)

	for _, tt := range data {
    b := bytes
		for _, want := range ints { var got uint32 
			b, got = tt.f(b)
			if got != want {
				t.Errorf("%s: %d(%x) != %d(%x)", tt.name, want, want, got, got)
			}
		}
	}
}

func BenchmarkLib(b *testing.B) {
  for i:=0; i<b.N; i++ {
    bytes := randBytes
    for _, want := range(randInts) {
      var got uint32
      bytes, got = lib(bytes)
      if got != want {
        b.Fatalf("%d(%x) != %d(%x)",got, got, want, want)
      }
    }
  }
}

func BenchmarkPlain(b *testing.B) {
  for i:=0; i<b.N; i++ {
    bytes := randBytes
    for _, want := range(randInts) {
      var got uint32
      bytes, got = plain(bytes)
      if got != want {
        b.Fatalf("%d(%x) != %d(%x)",got, got, want, want)
      }
    }
  }
}

func BenchmarkClean(b *testing.B) {
  for i:=0; i<b.N; i++ {
    bytes := randBytes
    for _, want := range(randInts) {
      var got uint32
      bytes, got = clean(bytes)
      if got != want {
        b.Fatalf("%d(%x) != %d(%x)",got, got, want, want)
      }
    }
  }
}

func BenchmarkMagic(b *testing.B) {
  for i:=0; i<b.N; i++ {
    bytes := randBytes
    for _, want := range(randInts) {
      var got uint32
      bytes, got = magic(bytes)
      if got != want {
        b.Fatalf("%d(%x) != %d(%x)",got, got, want, want)
      }
    }
  }
}
