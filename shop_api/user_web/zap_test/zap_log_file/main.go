package main

import (
	"go.uber.org/zap"
	"time"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig() //
	cfg.OutputPaths = []string{
		"stderr",
		"stdout",
		"./zap.log"} //输出到控制台和文件
	return cfg.Build()

}

func main() {
	logger, err := NewLogger() //生产环境
	if err != nil {
		panic(err)
	}
	//zap.NewDevelopment()             //开发环境
	defer logger.Sync() // flushes buffer, if any
	url := "https://www.baidu.com/"

	sugar := logger.Sugar() //可以简化使用
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)

}
