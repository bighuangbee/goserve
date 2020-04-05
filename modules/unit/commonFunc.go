package unit

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

func RandStr(n int) string {

	const letter = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letter_len := len(letter)

	b := make([]byte, n)
	for i := range b {
		b[i] = letter[rand.Intn(letter_len)]
	}
	return string(b)
}

func MD5(str string) string {
	m := md5.New()
	m.Write([]byte (str))
	return hex.EncodeToString(m.Sum(nil))
}

/*
 对切片的元素去重
 */
func SliceUnique(slice []string) []string {
	result := make([]string, 0, len(slice))
	temp := map[string]struct{}{}
	for _, item := range slice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func SubStr(str string, start, length int) string{
	rs := []rune(str)
	return string(rs[start:length])
}