package utils

import (
	"crypto/rand"
	"math/big"
	"net/url"
)

func GenerateRandomString(l int) string {
	var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, l)
	strLen := int64(len(symbols))
	for i := range b {
		index, err := rand.Int(rand.Reader, big.NewInt(strLen))
		if err != nil {
			panic(err)
		}
		b[i] = symbols[index.Int64()]
	}
	return string(b)
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
