// Package crypt 提供加密和解密功能
package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// AESEncrypt 使用AES算法加密数据
// 参数:
//   - data: 需要加密的原始数据
//   - key: 加密密钥
//   - iv: 初始化向量(用于CBC模式)
//   - mode: 加密模式(CBC或ECB)
//   - padding: 填充方式(PKCS5或PKCS7)
//
// 返回:
//   - 加密后的数据
//   - 错误信息(如果有)
func AESEncrypt(data, key, iv []byte, mode Mode, padding Padding) ([]byte, error) {
	// 创建AES密码块
	block, err := aes.NewCipher(key)

	// 根据填充方式对数据进行填充
	switch padding {
	case PKCS5:
		data = PKCS5Padding(data, block.BlockSize())
	case PKCS7:
		data = PKCS7Padding(data, block.BlockSize())
	default:
		return nil, errors.New("unsupport padding " + string(padding))
	}

	// 创建用于存储加密结果的缓冲区
	encrypted := make([]byte, len(data))
	if err != nil {
		println(err.Error())
		return nil, err
	}

	var encrypter cipher.BlockMode

	// 根据加密模式创建相应的加密器
	switch mode {
	case CBC:
		// CBC模式需要IV
		encrypter = cipher.NewCBCEncrypter(block, iv)
	case ECB:
		// ECB模式不需要IV
		encrypter = NewECBEncrypter(block)
	default:
		return nil, errors.New("unsupport mode " + string(mode))
	}

	// 执行加密操作
	encrypter.CryptBlocks(encrypted, data)
	return encrypted, nil
}

// AESDecrypt 使用AES算法解密数据
// 参数:
//   - src: 需要解密的加密数据
//   - key: 解密密钥
//   - iv: 初始化向量(用于CBC模式)
//   - mode: 解密模式(CBC或ECB)
//   - padding: 填充方式(PKCS5或PKCS7)
//
// 返回:
//   - 解密后的原始数据
//   - 错误信息(如果有)
func AESDecrypt(src, key, iv []byte, mode Mode, padding Padding) (data []byte, err error) {
	// 创建用于存储解密结果的缓冲区
	decrypted := make([]byte, len(src))

	// 创建AES密码块
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		println(err.Error())
		return nil, err
	}

	// 根据解密模式创建相应的解密器
	var decrypter cipher.BlockMode
	switch mode {
	case ECB:
		// ECB模式不需要IV
		decrypter = NewECBDecrypter(block)
	case CBC:
		// CBC模式需要IV
		decrypter = cipher.NewCBCDecrypter(block, iv)
	default:
		return nil, errors.New("unsupport mode " + string(mode))
	}

	// 执行解密操作
	decrypter.CryptBlocks(decrypted, src)

	// 根据填充方式去除填充
	switch padding {
	case PKCS5:
		return PKCS5UnPadding(decrypted), nil
	case PKCS7:
		return PKCS7UnPadding(decrypted), nil
	default:
		return nil, errors.New("unsupport padding " + string(padding))
	}
}
