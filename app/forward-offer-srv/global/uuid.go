package global

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	var err error
	snowflake.Epoch = time.Now().Unix()
	node, err = snowflake.NewNode(88)
	if err != nil {
		panic(err)
	}
}

func GetUUID() string {
	return fmt.Sprintf("OK%s", node.Generate())
}

func GetOkexOrderID() string {
	return fmt.Sprintf("OK%s", node.Generate())
}