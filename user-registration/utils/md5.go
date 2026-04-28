package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5EncodeWithSalt(data, salt string) string {
	return MD5Encode(data + salt)
}
