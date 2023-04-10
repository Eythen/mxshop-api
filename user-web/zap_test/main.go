package main

import (
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction() //生产环境
	//logger, _ := zap.NewDevelopment() //开发环境
	defer logger.Sync()
	url := "https://imooc.com"

	//性能高
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("nums", 3),
	)
	//使用反射机制
	//sugar := logger.Sugar()
	//sugar.Infow("failed to fetch URL",
	//	"url", url,
	//	"attempt", 3,
	//)
	//sugar.Infof("failed to fetch URL: %s", url)
}
