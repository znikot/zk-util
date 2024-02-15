package crypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"

	"math/big"
)

var (
	errDataToLarge     = errors.New("message too long for RSA public key size")
	errDataLen         = errors.New("data length error")
	errDataBroken      = errors.New("data broken, first byte is not zero")
	errKeyPairDismatch = errors.New("data is not encrypted by the private key")
	errDecryption      = errors.New("decryption error")
	errPublicKey       = errors.New("get public key error")
	errPrivateKey      = errors.New("get private key error")
)

const (
	// mode encrypt with rsa public key
	modePublicEncrypt = iota
	// mode decrypt with rsa public key
	modePublicDecrypt
	// mode encrypt with rsa private key
	modePrivateEncrypt
	// mode decrypt with rsa private key
	modePrivateDecrypt
)

// encrypt data with private key
func PrivateEncrypt(key *rsa.PrivateKey, src []byte) ([]byte, error) {
	return doFinal(src, modePrivateEncrypt, key)
}

// decrypt data with pem private key
func PrivateDecrypt(key *rsa.PrivateKey, src []byte) ([]byte, error) {
	// priKey, err := ParsePrivateKey([]byte(key))
	// if err != nil {
	// 	return nil, err
	// }
	return doFinal(src, modePrivateDecrypt, key)
}

// encrypt data with pem public key
func PublicEncrypt(key string, src []byte) ([]byte, error) {
	pubKey, err := ParsePublicKey([]byte(key))
	if err != nil {
		return nil, err
	}

	return doFinal(src, modePublicEncrypt, pubKey)
}

// decrypt data with pem public key
func PublicDecrypt(key *rsa.PublicKey, src []byte) ([]byte, error) {
	return doFinal(src, modePublicDecrypt, key)
}

// encrypt data with rsa public key
func pubKeyDecrypt(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	k := (pub.N.BitLen() + 7) / 8
	if k != len(data) {
		return nil, errDataLen
	}
	m := new(big.Int).SetBytes(data)
	if m.Cmp(pub.N) > 0 {
		return nil, errDataToLarge
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)
	d := leftPad(m.Bytes(), k)
	if d[0] != 0 {
		return nil, errDataBroken
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, errKeyPairDismatch
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}

// encrypt data with rsa private key
func priKeyEncrypt(rand io.Reader, priv *rsa.PrivateKey, hashed []byte) ([]byte, error) {
	tLen := len(hashed)
	k := (priv.N.BitLen() + 7) / 8
	if k < tLen+11 {
		return nil, errDataLen
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], hashed)
	m := new(big.Int).SetBytes(em)
	c, err := decrypt(rand, priv, m)
	if err != nil {
		return nil, err
	}
	copyWithLeftPad(em, c.Bytes())
	return em, nil
}

// 公钥加密或解密Reader
func pubKeyIO(pub *rsa.PublicKey, in io.Reader, out io.Writer, isEncrypt bool) error {
	k := (pub.N.BitLen() + 7) / 8
	if isEncrypt {
		k = k - 11
	}
	buf := make([]byte, k)
	var b []byte
	var err error
	size := 0
	for {
		size, err = in.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if size < k {
			b = buf[:size]
		} else {
			b = buf
		}
		if isEncrypt {
			b, err = rsa.EncryptPKCS1v15(rand.Reader, pub, b)
		} else {
			b, err = pubKeyDecrypt(pub, b)
		}
		if err != nil {
			return err
		}
		if _, err = out.Write(b); err != nil {
			return err
		}
	}
}

// 私钥加密或解密Reader
func priKeyIO(pri *rsa.PrivateKey, r io.Reader, w io.Writer, isEncrypt bool) error {
	k := (pri.N.BitLen() + 7) / 8
	if isEncrypt {
		k = k - 11
	}
	buf := make([]byte, k)
	var err error
	var b []byte
	size := 0
	for {
		size, err = r.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if size < k {
			b = buf[:size]
		} else {
			b = buf
		}
		if isEncrypt {
			b, err = priKeyEncrypt(rand.Reader, pri, b)
		} else {
			b, err = rsa.DecryptPKCS1v15(rand.Reader, pri, b)
		}

		if err != nil {
			return err
		}
		if _, err = w.Write(b); err != nil {
			return err
		}
	}
}

var pemStart = []byte("-----BEGIN ")

// 读取公钥
// 公钥可以没有如 -----BEGIN PUBLIC KEY-----的前缀后缀
func ParsePublicKey(in []byte) (*rsa.PublicKey, error) {
	var pubKeyBytes []byte
	if bytes.HasPrefix(in, pemStart) {
		block, _ := pem.Decode(in)
		if block == nil {
			return nil, errPublicKey
		}
		pubKeyBytes = block.Bytes
	} else {
		var err error
		pubKeyBytes, err = base64.StdEncoding.DecodeString(string(in))
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil, errPublicKey
		}
	}

	pub, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	} else {
		return pub.(*rsa.PublicKey), err
	}

}

// 解释单行的密钥信息
func ParsePemLine(in []byte, private bool) []byte {
	tmp := make([]byte, 0)
	if private {
		tmp = append(tmp, []byte("-----BEGIN RSA PRIVATE KEY-----\n")...)
	} else {
		tmp = append(tmp, []byte("-----BEGIN RSA PUBLIC KEY-----\n")...)
	}
	rest := in
	for len(rest) > 64 {
		tmp = append(tmp, rest[0:64]...)
		tmp = append(tmp, '\n')
		rest = rest[64:]
	}
	if len(rest) > 0 {
		tmp = append(tmp, rest...)
	}
	if private {
		tmp = append(tmp, []byte("\n-----END RSA PRIVATE KEY-----")...)
	} else {
		tmp = append(tmp, []byte("\n-----END RSA PUBLIC KEY-----")...)
	}
	return tmp
}

// 读取私钥
func ParsePrivateKey(in []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(in)
	if block == nil {
		return nil, errPrivateKey
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return pri, nil
	}
	pri2, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	} else {
		return pri2.(*rsa.PrivateKey), nil
	}
}

// 从crypto/rsa复制
var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)

// 从crypto/rsa复制
func decrypt(random io.Reader, priv *rsa.PrivateKey, c *big.Int) (m *big.Int, err error) {
	if c.Cmp(priv.N) > 0 {
		err = errDecryption
		return
	}
	var ir *big.Int
	if random != nil {
		var r *big.Int

		for {
			r, err = rand.Int(random, priv.N)
			if err != nil {
				return
			}
			if r.Cmp(bigZero) == 0 {
				r = bigOne
			}
			var ok bool
			ir, ok = modInverse(r, priv.N)
			if ok {
				break
			}
		}
		bigE := big.NewInt(int64(priv.E))
		rpowe := new(big.Int).Exp(r, bigE, priv.N)
		cCopy := new(big.Int).Set(c)
		cCopy.Mul(cCopy, rpowe)
		cCopy.Mod(cCopy, priv.N)
		c = cCopy
	}

	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}
	if ir != nil {
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}

	return
}

// 从crypto/rsa复制
func copyWithLeftPad(dest, src []byte) {
	numPaddingBytes := len(dest) - len(src)
	for i := 0; i < numPaddingBytes; i++ {
		dest[i] = 0
	}
	copy(dest[numPaddingBytes:], src)
}

// 从crypto/rsa复制
func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}

// 从crypto/rsa复制
func modInverse(a, n *big.Int) (ia *big.Int, ok bool) {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)
	g.GCD(x, y, a, n)
	if g.Cmp(bigOne) != 0 {
		return
	}
	if x.Cmp(bigOne) < 0 {
		x.Add(x, n)
	}
	return x, true
}

// 构建加密/解密io
func buildIO(in io.Reader, out io.Writer, mode int, key interface{}) error {
	switch mode {
	case modePublicEncrypt:
		return pubKeyIO(key.(*rsa.PublicKey), in, out, true)
	case modePublicDecrypt:
		return pubKeyIO(key.(*rsa.PublicKey), in, out, false)
	case modePrivateEncrypt:
		return priKeyIO(key.(*rsa.PrivateKey), in, out, true)
	case modePrivateDecrypt:
		return priKeyIO(key.(*rsa.PrivateKey), in, out, false)
	default:
		return errors.New("mode not found")
	}
}

// 执行加解密
func doFinal(in []byte, mode int, key interface{}) ([]byte, error) {
	out := bytes.NewBuffer(nil)
	err := buildIO(bytes.NewReader(in), out, mode, key)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(out)
}
