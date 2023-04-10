package initialize

import (
	"github.com/gin-gonic/gin"
	router2 "mxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	ApiGroup := Router.Group("/v1")
	router2.InitUserRouter(ApiGroup)

	return Router
}
