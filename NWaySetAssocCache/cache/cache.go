package cache

import (
	l "container/list"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"math"
	"sync"
	"unsafe"
)

const (
	LRU = 0
	MRU = 1
	CUSTOM = 2
)
type Cachestrategy = uint8

type val struct {
	_v interface{}
	lp unsafe.Pointer
}

type bcache struct {
	_map   []map[interface{}]interface{}
	_lck   []sync.Mutex
	cap    int
	limit  int
	sets   int
	_size  int
	_ktype string
	_vtype string
	stra   Cachestrategy
}

func (c *bcache) _checktype(k interface{}, v ...interface{}) bool{
	if c.empty() {
		return true
	}
	if (v == nil) {
		return c._ktype == fmt.Sprintf("%T", k)
	}
	return  c._ktype == fmt.Sprintf("%T", k) &&  c._vtype == fmt.Sprintf("%T", v[0])
}

func (c *bcache) _updateType(k interface{}, v ...interface{}) {
	c._ktype = fmt.Sprintf("%T", k)
	if v != nil {
		c._vtype = fmt.Sprintf("%T", v[0])
	}
}

func (c *bcache) _setsize(s int) {
	c._size = s
}

func (c *bcache) empty() bool{
	return c._size == 0
}

func (c *bcache) size() int{
	return c._size
}

func (c *bcache) _locate(i interface{}) (int, error) {
	bytes, err := json.Marshal(i)
	return int(crc32.ChecksumIEEE(bytes)) % c.sets, err
}


type Icache interface {
	Get(k interface{}) (interface{}, error)
	Put(k interface{}, v interface{}) error
	Contains(k interface{}) (bool, error)
	Size() int
	Remove(k interface{}) error
	Dump() string
}

func GetCache(cap int, item int, strategy Cachestrategy, custom_strategy ...Icache) (Icache, error){
	sets := int(math.Ceil(float64(cap/item)))
	var oper Icache
	switch strategy {
	case 0, 1:
		oper = &i_cache{
			lists: make([]l.List, sets, sets),
			bcache: bcache{
				_map: make([]map[interface{}]interface{}, sets),
				_lck: make([]sync.Mutex, sets, sets),
				cap: cap,
				limit: item,
				sets: sets,
				stra: strategy,
			},
		}
	case 2:
		if custom_strategy == nil {
			return nil, errors.New("custom strategy not provide")
		}
		oper = custom_strategy[0]
	default:
		return nil, errors.New("invalid strategy")
	}

	return oper, nil
}




