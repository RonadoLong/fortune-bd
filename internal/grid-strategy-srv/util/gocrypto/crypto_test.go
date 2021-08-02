package gocrypto

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAesDecrypt(t *testing.T) {
	//解密
	data := "d69802e16a3298b3e7a5222c77e3a9894e3f114106dd613cf986baf3e00d7d2d90afc4930fa02595f8c543effa4387b0"
	dataBytes, _ := hex.DecodeString(data)
	secretByte, _ := AesDecrypt(dataBytes)
	fmt.Println(string(secretByte))
}
