package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/goods-web/config"
	"mxshop-api/goods-web/proto"
)

var (
	ServerConfig   *config.ServerConfig = &config.ServerConfig{}
	Trans          ut.Translator
	GoodsSrvClient proto.GoodsClient
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
)
