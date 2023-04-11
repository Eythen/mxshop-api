package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
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

	port := global.ServerConfig.Port

	/*logger, _ := zap.NewProduction()
	defer logger.Sync()
	suger := logger.Sugar()*/
	/*
		1.zap.S();可以获取一个全局的sugar, 可以让我们自己设置一个全局的logger
		2.日志是分级别的 debug info warn error fetal（级别升高）
		3.S函数和L函数提供了一个全局的安全访问logger的途径
	*/
	zap.S().Infof("启动服务器，端口：%d", port)

	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}
}
