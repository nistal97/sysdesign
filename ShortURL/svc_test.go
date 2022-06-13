package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

import b "github.com/nistal97/sysdesign/ShortURL/bitly"

func TestInvalidStategy(t *testing.T) {
	digpref5 := "00000"
	for i := 1;i < 10;i ++ {
		require.True(t, string(*b.Convert2_62ruler(uint64(i))) == digpref5 + strconv.Itoa(i))
	}
	for i := 10;i < 36;i ++ {
		require.True(t, string(*b.Convert2_62ruler(uint64(i))) == digpref5 + fmt.Sprintf("%c", 'a' + i - 10))
	}
	for i := 36;i < 62;i ++ {
		require.True(t, string(*b.Convert2_62ruler(uint64(i))) == digpref5 + fmt.Sprintf("%c", 'A' + i - 36))
	}

}
