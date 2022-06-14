package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

import b "github.com/nistal97/sysdesign/ShortURL/bitly"

func TestConvert2_62ruler(t *testing.T) {
	v := []byte{49}
	for i := float64(1);i < math.Pow(62.0, 4.0);i ++ {
		require.True(t, bytes.Compare(*b.Convert2_62ruler(uint64(i)), v) == 0)
		v = plus(v)
	}
}

func plus(bytes []byte) []byte{
	return _plus(bytes, len(bytes) - 1, 1)
}

func _plus(bytes []byte, idx int, add int) []byte {
	if idx == -1 {
		if add == 1 {
			bytes = append([]byte{49}, bytes...)
		}
		return bytes
	}
	if add == 0 {
		return bytes
	}

	b := bytes[idx]
	my_add := 0

	//plus 1
	//if we hit 9,z,Z, plus 1, change to another section
	if b == 57 || b == 122 || b == 90 {
		if b == 57 {
			bytes[idx] = byte(97)
		} else if b == 122 {
			bytes[idx] = byte(65)
		} else {
			bytes[idx] = byte(48)
			my_add = 1
		}
	} else {
		//tail ele always plus 1
		bytes[idx] += byte(1)
	}
	return _plus(bytes, idx - 1, my_add)
}
