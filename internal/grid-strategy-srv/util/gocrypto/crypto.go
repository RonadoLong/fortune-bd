package gocrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/xml"
	"hash"

	"github.com/json-iterator/go"
)

var keyText = []byte("astaxie12798akljzmknm.ahkjkljl;k")

//EncodeJSON encodes JSON
func EncodeJSON(v interface{}) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
}

//DecodeJSON decodes JSON
func DecodeJSON(b []byte, v interface{}) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(b, v)
}

//EncodeXML encodes XML
func EncodeXML(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

//DecodeXML decodes XML
func DecodeXML(b []byte, v interface{}) error {
	return xml.Unmarshal(b, v)
}

//Base64Encode encodes bytes into base64 bytes
func Base64Encode(enc *base64.Encoding, b []byte) []byte {
	encoded := make([]byte, enc.EncodedLen(len(b)))
	enc.Encode(encoded, b)
	return encoded
}

//Base64Decode decodes base64 bytes
func Base64Decode(enc *base64.Encoding, b []byte) ([]byte, error) {
	decoded := make([]byte, enc.DecodedLen(len(b)))
	n, err := enc.Decode(decoded, b)
	return decoded[:n], err
}

//Base64EncodeJSON encodes JSON struct with BASE64
func Base64EncodeJSON(enc *base64.Encoding, v interface{}) ([]byte, error) {
	//convert to JSON bytes
	b, err := EncodeJSON(v)
	if err != nil {
		return nil, err
	}
	//BASE64 encode to JSON bytes
	return Base64Encode(enc, b), nil
}

//Base64DecodeJSON decodes BASE64 bytes into JSON struct
func Base64DecodeJSON(enc *base64.Encoding, encoded []byte, v interface{}) error {
	//BASE64 decode to JSON bytes
	b, err := Base64Decode(enc, encoded)
	if err != nil {
		return err
	}
	//convert from JSON bytes
	return DecodeJSON(b, v)
}

// func CalcHMAC(msg, secret []byte) []byte {
// 	return makeHMAC(sha256.New, msg, secret)
// }

//ConstTimeEqual compares two byte slice in constant time
func ConstTimeEqual(x, y []byte) bool {
	return subtle.ConstantTimeCompare(x, y) == 1
}

func hashSum(h hash.Hash, msg []byte) []byte {
	h.Write(msg)
	return h.Sum(nil)
}

//HMAC generates HMAC signature for message bytes, using specified algorithm
func HMAC(f func() hash.Hash, msg, secret []byte) []byte {
	return hashSum(hmac.New(f, secret), msg)
}

//SHA1 encodes message with SHA-1 algorithm
func SHA1(msg []byte) []byte {
	return hashSum(sha1.New(), msg)
}

//SHA256 encodes message with SHA-256 algorithm
func SHA256(msg []byte) []byte {
	return hashSum(sha256.New(), msg)
}

//MD5 encodes message with MD5 algorithm
func MD5(msg []byte) []byte {
	return hashSum(md5.New(), msg)
}

// AesEncrypt Aes加密 hex.EncodeToString()转字符串
func AesEncrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(keyText)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyText[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt 解密 string(b)转字符串
func AesDecrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(keyText)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyText[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// PKCS7Padding PKCS7Padding
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding PKCS7UnPadding
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// Base64EncodeUser Base64EncodeUser
func Base64EncodeUser(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

// Base64DecodeUser Base64DecodeUser
func Base64DecodeUser(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}
