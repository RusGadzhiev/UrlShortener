package encoder

const (
	alphabet    = "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	alphabetLen = len(alphabet)
	urlLen      = 10
)

var decoder map[byte]int

func init() {
	decoder = make(map[byte]int, alphabetLen)
	for i, v := range alphabet {
		decoder[byte(v)] = i
	}
}

func Encode(token int) string {
	shorten := [urlLen]byte{}
	for i := 0; i < urlLen; i++ {
		if token > 0 {
			shorten[urlLen-1-i] = alphabet[token%alphabetLen]
			token /= alphabetLen
		} else {
			shorten[urlLen-1-i] = alphabet[0]
		}
	}
	return string(shorten[:])
}