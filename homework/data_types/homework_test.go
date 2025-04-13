package main

import (
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

const maxSize = 8

var bufGlob = make([]byte, maxSize)

type uintegerable interface {
	uint16 | uint32 | uint64
}

func ToLittleEndian[T uintegerable](number T) T {
	size := unsafe.Sizeof(number)
	buf := bufGlob[:size]

	*(*T)(unsafe.Pointer(&buf[0])) = number

	s := int(size)
	for i := s/2 - 1; i >= 0; i-- {
		buf[i], buf[s-i-1] = buf[s-i-1], buf[i]
	}

	return *(*T)(unsafe.Pointer(&buf[0]))
}

func ToLittleEndianWithStd[T uintegerable](number T) T {
	switch v := any(number).(type) {
	case uint16:
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, v)
		return T(binary.LittleEndian.Uint16(bytes))
	case uint32:
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, v)
		return T(binary.LittleEndian.Uint32(bytes))
	case uint64:
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, v)
		return T(binary.LittleEndian.Uint64(bytes))
	default:
		return 0
	}
}

func TestСonversion(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)

			resultWithStd := ToLittleEndianWithStd(test.number)
			assert.Equal(t, test.result, resultWithStd)
		})
	}
}

func TestСonversion16(t *testing.T) {
	tests := map[string]struct {
		number uint16
		result uint16
	}{
		"test case #1": {
			number: 0x0000,
			result: 0x0000,
		},
		"test case #2": {
			number: 0xFFFF,
			result: 0xFFFF,
		},
		"test case #3": {
			number: 0x00FF,
			result: 0xFF00,
		},
		"test case #4": {
			number: 0xFF00,
			result: 0x00FF,
		},
		"test case #5": {
			number: 0x0102,
			result: 0x0201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)

			resultWithStd := ToLittleEndianWithStd(test.number)
			assert.Equal(t, test.result, resultWithStd)
		})
	}
}

func TestСonversion64(t *testing.T) {
	tests := map[string]struct {
		number uint64
		result uint64
	}{
		"test case #1": {
			number: 0x0000000000000000,
			result: 0x0000000000000000,
		},
		"test case #2": {
			number: 0xFFFFFFFFFFFFFFFF,
			result: 0xFFFFFFFFFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF00FF00FF,
			result: 0xFF00FF00FF00FF00,
		},
		"test case #4": {
			number: 0x00000000FFFFFFFF,
			result: 0xFFFFFFFF00000000,
		},
		"test case #5": {
			number: 0x0102030405060708,
			result: 0x0807060504030201,
		},
		"test case #6": {
			number: 0x01020304,
			result: 0x0403020100000000,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)

			resultWithStd := ToLittleEndianWithStd(test.number)
			assert.Equal(t, test.result, resultWithStd)
		})
	}
}

var Result16 uint16

func BenchmarkСonversion16(b *testing.B) {
	var value uint16 = 0x0102
	for i := 0; i < b.N; i++ {
		Result16 = ToLittleEndian(value)
	}
}

func BenchmarkСonversion16Std(b *testing.B) {
	var value uint16 = 0x0102
	for i := 0; i < b.N; i++ {
		Result16 = ToLittleEndianWithStd(value)
	}
}

var Result32 uint32

func BenchmarkСonversion32(b *testing.B) {
	var value uint32 = 0x01020304
	for i := 0; i < b.N; i++ {
		Result32 = ToLittleEndian(value)
	}
}

func BenchmarkСonversion32Std(b *testing.B) {
	var value uint32 = 0x01020304
	for i := 0; i < b.N; i++ {
		Result32 = ToLittleEndianWithStd(value)
	}
}

var Result64 uint64

func BenchmarkСonversion64(b *testing.B) {
	var value uint64 = 0x0102030405060708
	for i := 0; i < b.N; i++ {
		Result64 = ToLittleEndian(value)
	}
}

func BenchmarkСonversion64Std(b *testing.B) {
	var value uint64 = 0x0102030405060708
	for i := 0; i < b.N; i++ {
		Result64 = ToLittleEndianWithStd(value)
	}
}
