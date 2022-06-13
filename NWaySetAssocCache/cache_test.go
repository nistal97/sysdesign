package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)
import c "github.com/nistal97/sysdesign/NWaySetAssocCache/cache"

//test invalid cache strategy provided
func TestInvalidStategy(t *testing.T) {
	//invalid strategy choice.
	if _, err := c.GetCache(1024, 16, 127); err != nil {
		fmt.Printf("Expected Fail:%s\n", err)
	}
	//indicating custom cache strategy but not provided
	if _, err := c.GetCache(1024, 16, c.CUSTOM); err != nil {
		fmt.Printf("Expected Fail:%s\n", err)
	}
}

//test put,get,remove,lru,mru
func TestPutGetRemove(t *testing.T) {
	//rtype: key:0..26 val:a...z
	if cache, err := c.GetCache(32, 8, c.LRU); err != nil {
		fmt.Printf("Error occured:%s\n", err)
	} else {
		rt_putThenGet(t, cache)
	}
	if cache, err := c.GetCache(32, 8, c.MRU); err != nil {
		fmt.Printf("Error occured:%s\n", err)
	} else {
		rt_putThenGet(t, cache)
	}

	//uctype: key: struct{"a", 1} val: same:
	if cache, err := c.GetCache(32, 8, c.LRU); err != nil {
		fmt.Printf("Error occured:%s\n", err)
	} else {
		uct_putthenget(t, cache)
	}
	if cache, err := c.GetCache(32, 8, c.MRU); err != nil {
		fmt.Printf("Error occured:%s\n", err)
	} else {
		rt_putThenGet(t, cache)
	}

	//remove
	if cache, err := c.GetCache(1024, 64, c.LRU); err != nil {
		fmt.Printf("Error occured:%s\n", err)
	} else {
		for i := 0;i < 1024;i ++ {
			cache.Put(i, i)
		}
		for i := 0;i < 1024;i ++ {
			if i % 100 == 0 {
				err := cache.Remove(i)
				require.True(t, err == nil)
			}
		}
		for i := 0;i < 1024;i ++ {
			if i % 100 == 0 {
				b, _ := cache.Contains(i)
				require.True(t, !b)
			}
		}
	}

	//lru
	if cache, err := c.GetCache(4, 2, c.LRU); err != nil {
		fmt.Printf("Error occured:%s\n", err)
	} else {
		cache.Put("a", 1)
		v, _ := cache.Get("a")
		require.True(t, v == 1)
		cache.Put("b", 2)
		v, _ = cache.Get("b")
		require.True(t, v == 2)
		cache.Put("c", 3)
		v, _ = cache.Get("c")
		require.True(t, v == 3)
		cache.Put("d", 4)
		v, _ = cache.Get("d")
		require.True(t, v == 4)
		//get a so a should not be removed
		cache.Get("a")
		fmt.Println("lru before:", cache.Dump())
		cache.Put("e", 5)
		fmt.Println("lru after:", cache.Dump())
		require.True(t, cache.Size() == 4)
		b, _ := cache.Contains("a")
		require.True(t, b)
	}
	//mru
	if cache, err := c.GetCache(4, 2, c.MRU); err != nil {
		fmt.Printf("Error occured:%s\n", err)
	} else {
		cache.Put("a", 1)
		v, _ := cache.Get("a")
		require.True(t, v == 1)
		cache.Put("b", 2)
		v, _ = cache.Get("b")
		require.True(t, v == 2)
		cache.Put("c", 3)
		v, _ = cache.Get("c")
		require.True(t, v == 3)
		cache.Put("d", 4)
		v, _ = cache.Get("d")
		require.True(t, v == 4)
		//get a so a should  be removed
		cache.Get("a")
		fmt.Println("mru before:", cache.Dump())
		cache.Put("e", 5)
		fmt.Println("mru after:", cache.Dump())
		require.True(t, cache.Size() == 4)
		b, _ := cache.Contains("a")
		require.True(t, !b)
	}
}

func TestConcurrent(t *testing.T) {
	if cache, err := c.GetCache(4, 4, c.LRU); err != nil {
		fmt.Printf("Error occured:%s\n", err)
	} else {
		var wg sync.WaitGroup
		for i := 0;i < 1000;i ++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				cache.Put(i, i)
			}(i)
		}
		wg.Wait()
		require.True(t, cache.Size() == 4)
		fmt.Println(cache.Dump())
	}
}

type M struct {}
func (m M) Get(k interface{}) (interface{}, error) {
	return nil, errors.New("get failed")
}
func (m M)  Put (k interface{}, v interface{}) error {
return errors.New("put failed")
}
func (m M) Contains(k interface{}) (bool, error) {
	return false, errors.New("contains failed")
}
func (m M) Size () int {
	return 0
}
func (m M) Dump () string {
	return "customized cache"
}
func (m M) Remove(k interface{}) error {
	return errors.New("removed failed")
}
func TestCustomize(t *testing.T) {
	if cache, err := c.GetCache(4, 4, c.CUSTOM, M{}); err != nil {
		fmt.Printf("Error occured:%s\n", err)
	} else {
		err := cache.Put(1, 1)
		require.True(t, err != nil)
		fmt.Println(err)
		_, err = cache.Get(1)
		require.True(t, err != nil)
		fmt.Println(err)
		fmt.Println(cache.Dump())
	}
}

func rt_putThenGet(t *testing.T, cache c.Icache) {
	for i := 0;i <= 'z' - 'a';i ++ {
		require.True(t, cache.Put(i, fmt.Sprintf("%c", 'a' + i)) == nil)
	}
	for i := 0;i <= 'z' - 'a';i ++ {
		v, _ := cache.Get(i)
		require.True(t, v == fmt.Sprintf("%c", 'a' + i))
	}

	//try invalid type
	err := cache.Put("a", 1)
	require.True(t, err != nil)
	v, _ := cache.Get(float32(0))
	require.True(t, v == nil)
}

type kv struct {
	Fa string
	Fb int
}
func uct_putthenget(t *testing.T, cache c.Icache) {
	for i := 0;i <= 'z' - 'a';i ++ {
		kv := kv{ Fa: fmt.Sprintf("%c", 'a' + i), Fb: i}
		require.True(t, cache.Put(kv, kv) == nil)
	}
	for i := 0;i <= 'z' - 'a';i ++ {
		_kv := kv{ Fa: fmt.Sprintf("%c", 'a' + i), Fb: i}
		_v, _ := cache.Get(_kv)
		v := _v.(kv)
		require.True(t, v.Fa == fmt.Sprintf("%c", 'a' + i))
		require.True(t, v.Fb == i)
	}

	//try invalid type
	err := cache.Put("a", 1)
	require.True(t, err != nil)
	v, _ := cache.Get(float32(0))
	require.True(t, v == nil)
}





