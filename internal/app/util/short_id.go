package util

import "math/rand"

func GenerateShortId(length int) string {
	const base62String = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	shortId := make([]byte, length)
	for i := 0; i < length; i++ {
		shortId[i] = base62String[rand.Intn(len(base62String))]
	}

	return string(shortId)
}
