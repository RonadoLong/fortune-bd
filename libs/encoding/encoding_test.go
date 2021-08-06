package encoding

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	nihao := "删除"
	data, _ := AesEncrypt([]byte(nihao))
	log.Println(hex.EncodeToString(data))

}
