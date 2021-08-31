package ucode

import (
	"log"
	"testing"
)

func TestGetRandomString(t *testing.T) {

	got := GetRandomString(8)
	log.Println(got)
}

func TestRandStringRunes(t *testing.T) {

}
