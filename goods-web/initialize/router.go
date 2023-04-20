package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/middlewares"
	router2 "mxshop-api/goods-web/router"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	
	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/v1")
	router2.InitGoodsRouter(ApiGroup)

	return Router
}
