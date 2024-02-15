package crypt

import (
	"crypto/cipher"

	"github.com/tjfoc/gmsm/sm4"
)

// SM4EncryptBlocks 按块加密/解密
func SM4EncryptBlocks(key, src []byte) (result []byte, err error) {
	result = make([]byte, 0)

	idx, end := 0, sm4.BlockSize
	l := len(src)
	for idx < l {
		t := sm4.BlockSize
		if end >= l {
			end = l
			t = l - idx
		}
		// log.Infof("","len %d, index %d, end %d, t %d", l, idx, end, t)
		tmp := make([]byte, t)
		copy(tmp, src[idx:end])
		// dist := make([]byte, t)
		// sd.EncryptBlock(key, dist, tmp)
		var dist []byte
		dist, err = sm4.Sm4Ecb(key, tmp, true)
		if err != nil {
			return nil, err
		}
		result = append(result, dist[0:t]...)
		idx += sm4.BlockSize
		end += sm4.BlockSize
	}

	return
}

// SM4DecryptBlocks 按块解密
func SM4DecryptBlocks(key, src []byte) (result []byte, err error) {
	result = make([]byte, 0)

	idx, end := 0, sm4.BlockSize
	l := len(src)
	for idx < l {
		t := sm4.BlockSize
		if end >= l {
			end = l
			t = l - idx
		}
		tmp := make([]byte, t)
		copy(tmp, src[idx:end])
		// dist := make([]byte, t)
		// sd.DecryptBlock(key, dist, tmp)
		var dist []byte
		dist, err = sm4.Sm4Ecb(key, tmp, false)
		if err != nil {
			return nil, err
		}

		result = append(result, dist[0:t]...)
		end += sm4.BlockSize
		idx += sm4.BlockSize
	}

	return
}

// encrtyp data with SM4
func SM4Encrypt(plainText, key, keyiv []byte, mode Mode, padding Padding) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(keyiv) == 0 {
		keyiv = make([]byte, sm4.BlockSize)
	}
	blockSize := block.BlockSize()
	var origData []byte
	switch padding {
	case PKCS5:
		origData = PKCS5Padding(plainText, blockSize)
	case PKCS7:
		origData = PKCS7Padding(plainText, blockSize)
	case ZERO:
		origData = ZeroPadding(plainText, blockSize)
	case NONE:
		//
	default:
		origData = PKCS5Padding(plainText, blockSize)
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
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return cryted, nil
}

// decrypt data with SM4
func SM4Decrypt(cipherText, key, keyiv []byte, mode Mode, padding Padding) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(keyiv) == 0 {
		keyiv = make([]byte, sm4.BlockSize)
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
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
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
		origData = PKCS5UnPadding(origData)
	}
	return origData, nil
}
