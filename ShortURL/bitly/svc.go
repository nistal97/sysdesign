package bitly

type codec interface {
	Encode(url string)*[]byte
	Decode(url string)*[]byte
}

type ShortURLCodec interface {
	codec
}





