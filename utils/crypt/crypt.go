package crypt

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Guid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

func Crc32(str string) string {
	u32 := crc32.ChecksumIEEE([]byte(str))
	return fmt.Sprint(u32)
}

func C32() string {
	return Crc32(Guid())
}
