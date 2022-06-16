package bitly

import "sync"

const (
	NUM_MAP    = 0x00
	HASH_MAP   = 0x01
)

type BitlyCodec struct {
	Generator int
	daemon sync.Once
}

func (m *BitlyCodec) Encode(url string) *[]byte {
	m.daemon.Do(func() {

	})

	return nil
}

//A..Za..z0..9
func EncodeNum2Bytes(n uint64)*[]byte {
	s := []byte{}
	for {
		if n == 0  {
			break
		}
		mod := byte(n % 62)
		offset := byte(0)
		if mod < 10 {
			offset = 48 + mod
		} else if mod < 36 {
			offset = 97 + mod - 10
		} else {
			offset = 65 + mod - 36
		}
		s = append([]byte{offset}, s...)
		n /= 62
	}
	return &s
}

