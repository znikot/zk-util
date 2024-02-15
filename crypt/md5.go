package crypt

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// calculate md5 sum of string
func Md5(src string) string {
	ctx := md5.New()
	ctx.Write([]byte(src))
	return hex.EncodeToString(ctx.Sum(nil))
}

// calculate md5 sum of file
func Md5File(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
