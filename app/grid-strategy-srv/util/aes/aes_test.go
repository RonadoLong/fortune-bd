package aes

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	secretSrc := "123456"
	salt := []byte("0okm9ijn8uhb7ygv")

	fmt.Println("加密前:", secretSrc)
	secretText, err := Encrypt(secretSrc, salt)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("加密后:", string(secretText))

	secretSrc2, err := Decrypt(secretText, salt)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("解密后:", string(secretSrc2))

	if secretSrc != secretSrc2 {
		t.Error("加密前和解密后值不一致")
	}
}
