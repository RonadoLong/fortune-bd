package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOffset(c *gin.Context) (int, int) {

	pageNum := GetPage(c)
	pageSize := GetPageSize(c)

	offset := (pageNum - 1) * pageSize
	return offset, pageSize
}

func GetPage(c *gin.Context) int {
	ret, _ := strconv.Atoi(c.Param("pageNum"))
	if 1 > ret {
		ret = 1
	}
	return ret
}

func GetPageSize(c *gin.Context) int {
	ret, _ := strconv.Atoi(c.Param("pageSize"))
	if 1 > ret {
		ret = 15
	}
	return ret
}

func GeneratorPage(pageNum int, pageSize int) (int, int) {
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	pageNum = (pageNum - 1) * pageSize
	return pageNum, pageSize
}
