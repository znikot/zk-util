package crypt

import (
	"bytes"
)

// encrypt mode
type Mode string

// padding
type Padding string

const (
	//cbc
	CBC Mode = "CBC"
	//ecb
	ECB = "ECB"

	//PKCS5 padding
	PKCS5 Padding = "PKCS5"
	//PKCS7 padding
	PKCS7 = "PKCS7"
	//ZERO 0 padding
	ZERO = "ZERO"
	//NONE padding
	NONE = "NONE"
)

// PKCS5Padding
func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// PKCS5UnPadding
func PKCS5UnPadding(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

// PKCS7Padding
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// ZeroPadding 补齐0
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

// ZeroUnPadding 去0
func ZeroUnPadding(origData []byte) []byte {
	// result := make([]byte, 0)
	// for _, b := range origData {
	// 	// log.Infof("","%d", b)
	// 	if b != 0 {
	// 		result = append(result, b)
	// 	}
	// }
	// return result
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}
