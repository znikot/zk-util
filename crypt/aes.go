package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// encrypt data with aes
func AESEncrypt(data, key []byte, mode Mode, padding Padding) ([]byte, error) {
	block, err := aes.NewCipher(key)

	//确定对齐模式
	switch padding {
	case PKCS5:
		data = PKCS5Padding(data, block.BlockSize())
	case PKCS7:
		data = PKCS7Padding(data, block.BlockSize())
	default:
		return nil, errors.New("unsupport padding " + string(padding))
	}
	encrypted := make([]byte, len(data))
	if err != nil {
		println(err.Error())
		return nil, err
	}
	var encrypter cipher.BlockMode

	//获取加密模式
	switch mode {
	case CBC:
		encrypter = cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	case ECB:
		encrypter = NewECBEncrypter(block)
	default:
		return nil, errors.New("unsupport mode " + string(mode))
	}

	encrypter.CryptBlocks(encrypted, data)
	return encrypted, nil
}

// decrypt data with aes
func AESDecrypt(src, key []byte, mode Mode, padding Padding) (data []byte, err error) {
	// log.Debugf("AESDecrypt using key: %x", key)
	decrypted := make([]byte, len(src))
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	var decrypter cipher.BlockMode
	switch mode {
	case ECB:
		decrypter = NewECBDecrypter(block)
	case CBC:
		decrypter = cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	default:
		return nil, errors.New("unsupport mode " + string(mode))
	}
	decrypter.CryptBlocks(decrypted, src)

	switch padding {
	case PKCS5:
		return PKCS5UnPadding(decrypted), nil
	case PKCS7:
		return PKCS7UnPadding(decrypted), nil
	default:
		return nil, errors.New("unsupport padding " + string(padding))
	}
}
