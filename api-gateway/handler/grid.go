package handler

import (
	"github.com/gin-gonic/gin"
	grid_strategy_srv "wq-fotune-backend/service/grid-strategy-srv"
)

func InitGridEngine(engine *gin.RouterGroup) {
	go grid_strategy_srv.InitMain(engine)
}
