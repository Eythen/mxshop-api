package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/initialize"
	"mxshop-api/goods-web/utils"
	"mxshop-api/goods-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"
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

	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := registerClient.Register(global.ServerConfig.Host, port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}

	zap.S().Infof("启动服务器，端口：%d", port)

	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
			zap.S().Panic("启动失败：", err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Info("注销服务失败：", err.Error())
	} else {
		zap.S().Info("注销服务成功")
	}
}
