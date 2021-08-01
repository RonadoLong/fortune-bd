package helper

import (
	"fmt"
	"testing"
)

func TestRepalceSpan(t *testing.T) {
	var id = "Spam\n4f4b8856-c960-45fe-a4e5-5f5942cd017b"
	span := TrimSpan(id)
	fmt.Println(span)
}
