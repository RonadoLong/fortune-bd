package snowflake

import (
	"fortune-bd/libs/logger"
	"github.com/bwmarrin/snowflake"
	"os"
)

var SNode *snowflake.Node

func init() {
	snode, err := snowflake.NewNode(1)
	if err != nil {
		logger.Errorf("初始化 snowflake 失败 %v", err)
		os.Exit(-1)
	}
	SNode = snode
}
