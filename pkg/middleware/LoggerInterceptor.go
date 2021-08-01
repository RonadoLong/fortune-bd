package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"wq-fotune-backend/libs/logger"
)

const limitSize = 300

type bodyLogWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func InputOutputLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		//  处理前打印输入信息
		buf := bytes.Buffer{}
		_, _ = buf.ReadFrom(c.Request.Body)
		method := c.Request.Method
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			logger.Info("<<<<<<",
				logger.String("method", c.Request.Method),
				logger.Any("url", c.Request.URL),
				logger.Int("size", buf.Len()),
				logger.String("body", getBodyData(buf)),
			)
		} else {
			logger.Info("<<<<<<",
				logger.String("method", c.Request.Method),
				logger.Any("url", c.Request.URL),
			)
		}
		c.Request.Body = ioutil.NopCloser(&buf)

		//  替换writer
		newWriter := &bodyLogWriter{body: bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = newWriter

		// 处理请求
		c.Next()

		// 处理后打印返回信息
		logger.Info(">>>>>>",
			logger.Int("code", c.Writer.Status()),
			logger.String("method", c.Request.Method),
			logger.String("url", c.Request.URL.Path),
			logger.String("time", fmt.Sprintf("%dus", time.Now().Sub(start).Nanoseconds()/1000)),
			logger.Int("size", newWriter.body.Len()),
			logger.String("response", strings.TrimRight(getBodyData(newWriter.body), "\n")),
		)
	}
}

func getBodyData(buf bytes.Buffer) string {
	var body string
	if buf.Len() > limitSize {
		body = string(buf.Bytes()[:limitSize]) + " ...... "
	} else {
		body = buf.String()
	}
	return body
}
