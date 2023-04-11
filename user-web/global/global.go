package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/user-web/config"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans        ut.Translator
)
