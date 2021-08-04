package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 填充数据
func paddingData(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func removePaddingData(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

// Encrypt 加密
func Encrypt(secretSrc string, salt []byte) (string, error) {
	src := []byte(secretSrc)

	block, err := aes.NewCipher(salt)
	if err != nil {
		return "", err
	}

	src = paddingData(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, salt)
	blockMode.CryptBlocks(src, src)
	val := base64.StdEncoding.EncodeToString(src)

	return val, nil
}

// Decrypt 解密，secretText为base64编码后的数据
func Decrypt(secretText string, salt []byte) (string, error) {
	src, err := base64.StdEncoding.DecodeString(secretText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(salt)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, salt)
	blockMode.CryptBlocks(src, src)
	src = removePaddingData(src)

	return string(src), nil
}
