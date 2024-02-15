package crypt

import (
	"crypto/cipher"
	"crypto/des"
)

// encrypt data with des
func DESEncrypt(origData, key []byte, mode Mode, padding Padding) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	switch padding {
	case PKCS5:
		origData = PKCS5Padding(origData, block.BlockSize())
	case PKCS7:
		origData = PKCS7Padding(origData, block.BlockSize())
	case ZERO:
		origData = ZeroPadding(origData, block.BlockSize())
	case NONE:
		//
	default:
		origData = ZeroPadding(origData, block.BlockSize())
	}
	var blockMode cipher.BlockMode
	switch mode {
	case CBC:
		blockMode = cipher.NewCBCEncrypter(block, key)
	case ECB:
		blockMode = NewECBEncrypter(block)
	default:
		blockMode = cipher.NewCBCEncrypter(block, key)
	}

	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// descrypt data with des
func DESDecrypt(crypted, key []byte, mode Mode, padding Padding) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var blockMode cipher.BlockMode
	switch mode {
	case CBC:
		blockMode = cipher.NewCBCDecrypter(block, key)
	case ECB:
		blockMode = NewECBDecrypter(block)
	default:
		blockMode = cipher.NewCBCDecrypter(block, key)
	}
	origData := make([]byte, len(crypted))

	blockMode.CryptBlocks(origData, crypted)
	switch padding {
	case PKCS5:
		origData = PKCS5UnPadding(origData)
	case ZERO:
		origData = ZeroUnPadding(origData)
	case NONE:
		//
	default:
		origData = ZeroUnPadding(origData)
	}
	return origData, nil
}

// encrypt data with 3des
func DES3Encrypt(origData, key, keyiv []byte, mode Mode, padding Padding) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	switch padding {
	case PKCS5:
		origData = PKCS5Padding(origData, block.BlockSize())
	case PKCS7:
		origData = PKCS7Padding(origData, block.BlockSize())
	case ZERO:
		origData = ZeroPadding(origData, block.BlockSize())
	case NONE:
		//
	default:
		origData = ZeroPadding(origData, block.BlockSize())
	}
	var blockMode cipher.BlockMode
	switch mode {
	case CBC:
		blockMode = cipher.NewCBCEncrypter(block, keyiv)
	case ECB:
		blockMode = NewECBEncrypter(block)
	default:
		blockMode = cipher.NewCBCEncrypter(block, keyiv)
	}
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// descrypt data with 3des
func DES3Decrypt(crypted, key, keyiv []byte, mode Mode, padding Padding) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	var blockMode cipher.BlockMode
	switch mode {
	case CBC:
		blockMode = cipher.NewCBCDecrypter(block, keyiv)
	case ECB:
		blockMode = NewECBDecrypter(block)
	default:
		blockMode = cipher.NewCBCDecrypter(block, keyiv)
	}
	origData := make([]byte, len(crypted))

	blockMode.CryptBlocks(origData, crypted)

	switch padding {
	case PKCS5:
		origData = PKCS5UnPadding(origData)
	case PKCS7:
		origData = PKCS7UnPadding(origData)
	case ZERO:
		origData = ZeroUnPadding(origData)
	case NONE:
		//
	default:
		origData = ZeroUnPadding(origData)
	}

	return origData, nil
}
