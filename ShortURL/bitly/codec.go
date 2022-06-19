package bitly

import (
	"hash/crc32"
	"sync"
	"sync/atomic"
)

type BitlyCodec struct {
	Generator int
	daemon sync.Once
	buckets []bucket
	counter uint64
}

type bucket struct {
	lck sync.Mutex
	hosts map[string]string
}


func (m *BitlyCodec) Encode(url string) string {
	atomic.AddUint64(&m.counter, 1)
	_k := EncodeNum2Bytes(m.counter)
	k := DOMAIN + string(*_k)

	m.buckets[m.locate(k)].lck.Lock()
	defer m.buckets[m.locate(k)].lck.Unlock()
	if m.buckets[m.locate(k)].hosts == nil {
		m.buckets[m.locate(k)].hosts = make(map[string]string)
	}

	m.buckets[m.locate(k)].hosts[k] = url
	return k
}

func (m *BitlyCodec) Decode(url string) string {
	m.buckets[m.locate(url)].lck.Lock()
	defer m.buckets[m.locate(url)].lck.Unlock()
	if bs, ok := m.buckets[m.locate(url)].hosts[url]; ok {
		return bs
	}
	return ""
}

func (m *BitlyCodec) locate(url string) int {
	return int(crc32.ChecksumIEEE([]byte(url))) % BUCKET_SIZE
}

//A..Za..z0..9
func EncodeNum2Bytes(n uint64)*[]byte {
	s := []byte{}
	for {
		if n == 0  {
			break
		}
		mod := byte(n % BUCKET_SIZE)
		offset := byte(0)
		if mod < 10 {
			offset = 48 + mod
		} else if mod < 36 {
			offset = 97 + mod - 10
		} else {
			offset = 65 + mod - 36
		}
		s = append([]byte{offset}, s...)
		n /= BUCKET_SIZE
	}
	return &s
}

