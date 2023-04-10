package main

import (
	"go.uber.org/zap"
	"time"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./myproject.log",
		"stderr",
		"stdout",
	}
	return cfg.Build()
}

func main() {
	//logger, _ := zap.NewProduction()
	logger, err := NewLogger()
	if err != nil {
		panic(err)
	}
	su := logger.Sugar()
	defer su.Sync()
	url := "https://imooc.com"
	su.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("nums", 3),
		zap.Duration("backoff", time.Second),
	)
}
