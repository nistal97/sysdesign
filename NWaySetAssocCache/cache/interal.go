package cache

import (
	l "container/list"
	"errors"
	"fmt"
	"unsafe"
)

type i_cache struct {
    lists []l.List
	bcache
}

func (c *i_cache) Get(k interface{}) (interface{}, error) {
	if idx, err := c._locate(k); err != nil {
		return nil, err
	} else {
		c.bcache._lck[idx].Lock()
		defer c.bcache._lck[idx].Unlock()
		if v, ok := c._map[idx][k]; ok {
			removed := c.lists[idx].Remove((*l.Element)(v.(val).lp))
			c.lists[idx].PushFront(removed)
			c._map[idx][k] = val{
				_v: v.(val)._v,
				lp: unsafe.Pointer(c.lists[idx].Front()),
			}
			return v.(val)._v, nil
		}
		return nil, nil
	}
}

func (c *i_cache) Put(k interface{}, v interface{}) error {
	if idx, err := c._locate(k); err != nil {
		return err
	} else {
		c.bcache._lck[idx].Lock()
		defer c.bcache._lck[idx].Unlock()
		if !c.bcache._checktype(k, v) {
			return errors.New("invalid type provided")
		}

		if c.lists[idx].Len() == c.bcache.limit {
			var old_k *l.Element
			if c.bcache.stra == LRU {
				old_k = c.lists[idx].Back()
			} else if c.bcache.stra == MRU {
				old_k = c.lists[idx].Front()
			}
			c.lists[idx].Remove(old_k)
			delete(c._map[idx], old_k.Value)
			c.bcache._setsize(c.bcache._size - 1)
		} else if (c.bcache.empty()) {
			c._updateType(k, v)
		}

		if c._map[idx] == nil {
			c._map[idx] = make(map[interface{}]interface{}, c.bcache.cap)
		}
		c.lists[idx].PushFront(k)
		c._map[idx][k] = val{
			_v: v,
			lp: unsafe.Pointer(c.lists[idx].Front()),
		}
		c.bcache._setsize(c.bcache._size + 1)
	}
	return nil

}

func (c *i_cache) Contains(k interface{}) (bool,error) {
	if idx, err := c._locate(k); err != nil {
		return false, err
	} else {
		c.bcache._lck[idx].Lock()
		defer c.bcache._lck[idx].Unlock()
		if _, ok := c._map[idx][k]; ok {
			return true, nil
		}
		return false, nil
	}
}

func (c *i_cache) Remove(k interface{}) error {
	if idx, err := c._locate(k); err != nil {
		return err
	} else {
		c.bcache._lck[idx].Lock()
		defer c.bcache._lck[idx].Unlock()
		if v, ok := c._map[idx][k]; ok {
			c.lists[idx].Remove((*l.Element)(v.(val).lp))
			delete(c._map[idx], k)
			c.bcache._setsize(c.bcache._size - 1)
		}
		return nil
	}
}

func (c *i_cache) Size() int {
	return c.bcache._size
}

func (c *i_cache) Dump() string {
	s := ""
	sets := 0
	num := 0
	for n, l := range c.lists {
		s += "["
		for i := l.Front(); i != nil; i = i.Next() {
			s += fmt.Sprintf("%s%d%s%v ", "Set:", n, " key:", i.Value)
			num ++
		}
		s += "] "
		sets ++
	}
	s += fmt.Sprintf("%s%d%s%d ", " total Sets:", sets, " cached:", num)
	return s
}
