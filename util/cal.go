package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// CalculateHash calculate hash
func CalculateHash(key, data string) string {
	h := sha256.New()
	h.Write([]byte(key + data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// CalculateMD5 calculate md5
func CalculateMD5(key, data string) string {
	h := md5.New()
	h.Write([]byte(key + data))
	return hex.EncodeToString(h.Sum(nil))
}
