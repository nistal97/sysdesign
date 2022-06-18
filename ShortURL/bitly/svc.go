package bitly

const (
	BUCKET_SIZE = 62
	DOMAIN = "https://bitly.com/"
)

type ShortURLCodec interface {
	Encode(url string)string
	Decode(url string)string
}

func GetShortURLCodec() ShortURLCodec {
	c := &BitlyCodec{
		buckets: make([]bucket, BUCKET_SIZE, BUCKET_SIZE),
	}

	return c
}





