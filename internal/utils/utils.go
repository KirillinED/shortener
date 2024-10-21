package utils

import (
	"hash/crc32"
)

type URL string

func (u URL) String() string {
	return string(u)
}

func (u URL) Short() string {
	hash := crc32.ChecksumIEEE([]byte(u.String()))

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
