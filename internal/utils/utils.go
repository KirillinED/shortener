package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"time"
)

func ShortURL(url string) string {
	hash := crc32.ChecksumIEEE([]byte(url))

	return Base62Encode(hash)
}

func Base62Encode(num uint32) string {
	const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""
	for num > 0 {
		remainder := num % 62
		result = string(base62Chars[remainder]) + result
		num /= 62
	}
	return result
}

func GenerateUserId() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x", time.Now().UnixNano(), binary.BigEndian.Uint64(b)), nil
}

func GenerateHashKey() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
