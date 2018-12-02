package handler

import (
	"github.com/gin-gonic/gin"
	"shop-micro/helper"
	"time"
)

func ClientEngine() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	AuthMiddleware := &helper.GinJWTMiddleware{
		Realm:            "test zone",
		SigningAlgorithm: "",
		Key:              []byte("secret key"),
		Timeout:          time.Hour * 24 * 3000,
		MaxRefresh:       time.Hour,
		Authenticator:    Login,
		Authorizator: func(userId string, c *gin.Context) bool {
			return true
		},
		PayloadFunc: nil,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		IdentityHandler:       nil,
		TokenLookup:           "header:Authorization",
		TokenHeadName:         "Bearer",
		TimeFunc:              time.Now,
		HTTPStatusMessageFunc: nil,
		PrivKey:               helper.GetRsaPriKey(),
		PubKey:                helper.GetRsaPublicKey(),
	}

	api := router.Group("/api/client")

	/**  global interceptor */
	api.Use(AuthMiddleware.MiddlewareParseUser)

	/**  user logic */
	userGroup := api.Group("/user")
	userGroup.POST("/login", AuthMiddleware.LoginHandler)
	userGroup.GET("/getCode/:phone",GetPhoneCode)

	auth := userGroup.Group("/auth")
	auth.Use(AuthMiddleware.MiddlewareFunc())
	//{
	//	auth.POST("/update", controller.UserRest.UpdateUserInfo)
	//	auth.GET("/userInfo", controller.UserRest.GetUserInfo)
	//	auth.GET("/saveAgreement", controller.UserRest.SaveAgreement)
	//	auth.GET("/getUserAgreement", controller.UserRest.GetUserAgreement)
	//	auth.GET("/refresh_token", AuthMiddleware.RefreshHandler)
	//	auth.GET("/getUserAddress", controller.UserRest.GetUserAddress)
	//	auth.POST("/saveUserAddress", controller.UserRest.SaveUserAddress)
	//	auth.GET("/getUserIntegralFlow/:pageNum/:pageSize", controller.UserRest.GetUserIntegralFlow)
	//}

	/** home logic */
	homeGroup := api.Group("/home")
	homeGroup.GET("/headers", FindHomeHeadList)
	//homeGroup.GET("/list", home.HomeRest.FindHomeList)

	/** news logic  */
	news := api.Group("/news")
	news.GET("/category/list", GetNewsCategoryList)
	news.GET("/list/:category/:pageNum/:pageSize", GetNewsList)
	news.GET("/homeList/:pageNum/:pageSize", GetNewsList)
	//news.GET("/detail/:newsId", News.NewsRest.FindNewsDetail)
	//news.GET("/like/:newsId", News.NewsRest.AddLikeById)

	/** video logic  */
	video := api.Group("/video")
	video.GET("/list/:category/:pageNum/:pageSize", FindVideoList)
	video.GET("/homeList/:pageNum/:pageSize", FindVideoList)
	//video.GET("/detail/:videoId", News.VideoRest.FindVideoDetail)
	video.GET("categoryList", FindVideoCategoryList)

	/** product logic  */
	//product := api.Group("/product")
	//product.GET("/nav/list", FindGoodsNavList)
	//product.GET("/list/:category/:pageNum/:pageSize", FindGoodsList)
	//product.GET("/homeList/:pageNum/:pageSize", Goods.GoodsRest.FindGoodsAllList)
	//product.GET("/detail/:goodsId", Goods.GoodsRest.FindGoodsDetail)

	/** cart logic */
	//cart := api.Group("/cart")
	//cart.Use(AuthMiddleware.MiddlewareFunc())
	//cart.POST("/confirm", Goods.ShoppingCartRest.AddCart)
	//cart.GET("/list", Goods.ShoppingCartRest.FindCartList)
	//cart.DELETE("/del/:cartId", Goods.ShoppingCartRest.DelCart)
	//
	///** order logic */
	//order := api.Group("/order")
	//order.Use(AuthMiddleware.MiddlewareFunc())
	//order.POST("/create", Order.OrderInfoRest.CreateOrder)
	//order.POST("/list", Order.OrderInfoRest.FindOrderInfoList)
	//order.PUT("/cancel/:orderId", Order.OrderInfoRest.CancelOrder)
	//
	//other := api.Group("/other")
	//other.POST("/callback", Order.OrderInfoRest.ReceiveIPN)

	/** logic logic */
	//service := api.Group("/service")
	//service.GET("/category/list", Service.ServiceRest.FindServiceCategoryList)
	//service.GET("/findServicePaymentList", Service.ServiceRest.FindServicePaymentList)
	//service.GET("/area/list/:areaname", Service.ServiceRest.FindAreaListByParentId)
	//service.POST("/room/list/:pageNum/:pageSize", Service.ServiceRest.FindServiceRoomByCategory)
	//
	//serviceAuth := service.Use(AuthMiddleware.MiddlewareFunc())
	//serviceAuth.POST("/auth/save", Service.ServiceRest.AddServiceRoom)
	//serviceAuth.POST("/job/auth/save", Service.ServiceRest.SaveServiceJob)
	//serviceAuth.GET("/auth/self/:status/:pageNum/:pageSize", Service.ServiceRest.FindSelfService)

	return router
}
