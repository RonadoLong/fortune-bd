package client

import (
	"context"
	"log"
	"testing"
)

func TestHomeClient_FindHomeHeadList(t *testing.T) {
	client := NewHomeClient()
	headList, err := client.FindHomeHeadList(context.TODO())
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(headList)
}