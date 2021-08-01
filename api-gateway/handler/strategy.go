package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
)
//func InitExOrderEngine(engine *gin.RouterGroup) {
//	exOrderService = client.NewExOrderClient()
//	orderService = client.NewForwardOfferClient()
//	group := engine.Group("/strategy")
//	r.POST(urlPattern, ReverseProxy())
//	r.GET(urlPattern, ReverseProxy())
//}
func ReverseProxy() gin.HandlerFunc {
	target := "localhost:9531"
	return func(c *gin.Context) {
		director := func(req *http.Request) {
			r := c.Request
			req = r
			req.URL.Scheme = "http"
			req.URL.Host = target
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
