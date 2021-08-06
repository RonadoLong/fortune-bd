package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Code int32       `json:"code"` // return code, 0 for succ
	Msg  string      `json:"msg"`  // message
	Data interface{} `json:"data"` // data object
}

func NewResultInternalErr(msg string) *Result {
	return &Result{
		Code: ERROR_INTERNAL_SERVER,
		Msg:  msg,
		Data: nil,
	}
}

func NewResultSuccess(data interface{}) *Result {
	return &Result{
		Code: SUCCESS_CODE,
		Msg:  "ok",
		Data: data,
	}
}

const (
	SUCCESS_CODE           = 0
	ERROR_CODE_NOT_FOUND   = 1404
	ERROR_CODE_BIND_JSON   = 1400
	ERROR_INTERNAL_SERVER  = 1500
	ERROR_CODE_WRONG_PARAM = 1400
	ERROR_CODE_Max_Req     = 1503
	ERROR_CODE_CREATE      = 1500
	ERROR_CODE_UPDATE      = 1500
	ERROR_CODE_DELETE      = 1500
)

const (
	BIND_JSON_ERROR = "传参格式有误！"
	INTERNAL_SERVER = "内部服务错误！"
	OK              = "ok"
)

// NewResult creates a result with Code=0, Msg="", Data=nil.
func NewErrorParam(c *gin.Context, msg string, data interface{}) {
	result := &Result{
		Code: ERROR_CODE_WRONG_PARAM,
		Msg:  msg,
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func NewErrorCreate(c *gin.Context, msg string, data interface{}) {
	result := &Result{
		Code: ERROR_CODE_CREATE,
		Msg:  msg,
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func NewErrorUpdate(c *gin.Context, msg string, data interface{}) {
	result := &Result{
		Code: ERROR_CODE_UPDATE,
		Msg:  msg,
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func NewErrorMaxReq(c *gin.Context, msg string, data interface{}) {
	result := &Result{
		Code: ERROR_CODE_Max_Req,
		Msg:  msg,
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func NewSuccess(c *gin.Context, data interface{}) {
	result := &Result{
		Code: SUCCESS_CODE,
		Msg:  OK,
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func NewNotFoundErr(c *gin.Context, msg string, data interface{}) {
	result := &Result{
		Code: ERROR_CODE_NOT_FOUND,
		Msg:  msg,
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func NewBindJsonErr(c *gin.Context, data interface{}) {
	result := &Result{
		Code: ERROR_CODE_BIND_JSON,
		Msg:  BIND_JSON_ERROR,
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func NewInternalServerErr(c *gin.Context, data interface{}) {
	result := &Result{
		Code: ERROR_INTERNAL_SERVER,
		Msg:  INTERNAL_SERVER,
		Data: data,
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}

func NewInternalServerWithMsgErr(c *gin.Context, msg string) {
	NewErrWithCodeAndMsg(c, ERROR_INTERNAL_SERVER, msg)
}

func NewErrWithCodeAndMsg(c *gin.Context, code int32, msg string) {
	result := &Result{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.AbortWithStatusJSON(http.StatusOK, result)
}
