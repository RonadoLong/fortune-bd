
package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Result.
type Result struct {
	Code int         `json:"code"` // return code, 0 for succ
	Msg  string      `json:"msg"`  // message
	Data interface{} `json:"data"` // data object
}

func NewResult() *Result {
	return &Result{
		Code: 0,
		Msg:  "",
		Data: nil,
	}
}

func CreateSuccess(c *gin.Context, data interface{}){
	json := NewResult()
	json.Data = data
	json.Msg = "sccess"
	json.Code = 1000

	c.JSON(
		http.StatusOK,
		json,
	)
}

func CreateNotContent(c *gin.Context) {
	json := NewResult()
	json.Data = nil
	json.Msg = "No More Content"
	json.Code = 1004

	c.JSON(
		http.StatusOK,
		json,
	)
}

func CreateError(c *gin.Context){
	json := NewResult()
	json.Data = nil
	json.Msg = "logic fail"
	json.Code = 400

	c.JSON(
		http.StatusBadRequest,
		json,
	)
}


func CreateErrorWithMsg(c *gin.Context, msg string){
	json := NewResult()
	json.Data = nil
	json.Msg = msg
	json.Code = 400

	c.JSON(
		http.StatusBadRequest,
		json,
	)
}

func CreateErrorParams(c *gin.Context){
	json := NewResult()
	json.Data = nil
	json.Msg = "error params"
	json.Code = 400

	c.JSON(
		http.StatusBadRequest,
		json,
	)
}

func CreateSuccessByList(c *gin.Context, total interface{}, content interface{}){
	json := NewResult()

	json.Data = gin.H{
		"total": total,
		"content": content,
	}
	json.Msg = "success"
	json.Code = 1000

	c.JSON(
		http.StatusOK,
		json,
	)
}


func CreateErrorRequest(c *gin.Context){
	json := NewResult()
	json.Data = nil
	json.Msg = "The request is frequent"
	json.Code = 400

	c.JSON(
		http.StatusForbidden,
		json,
	)
}
