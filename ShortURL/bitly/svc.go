package bitly

type auto_increment struct {

}

//A..Za..z0..9
func Convert2_62ruler(n uint64)*[]byte {
	s := []byte{48,48,48,48,48,48}
	i := 5
	for {
		if n == 0 || i < 0 {
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
		s[i] = offset
		n /= 62
		i --
	}
	return &s
}