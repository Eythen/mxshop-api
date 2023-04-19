package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/initialize"
	"mxshop-api/goods-web/utils"
)

func main() {
	//1.初始化logger
	initialize.InitLogger()

	//2.初始化配置文件
	initialize.InitConfig()

	//3.初始化routers
	Router := initialize.Routers()

	//4.初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	//5.初始化srv的连接
	initialize.InitSrvConn()

	viper.AutomaticEnv()
	//如果是本地开发环境端口号固定，线上环境启动获取端口号
	debug := viper.GetBool("PROG27B48B2C051")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	port := global.ServerConfig.Port

	zap.S().Infof("启动服务器，端口：%d", port)

	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
