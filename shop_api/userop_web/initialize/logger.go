package initialize

import "go.uber.org/zap"

func InitLogger() {
	//设置全局日志级别
	/*
		1.	zap.S()可以获取到全局的logger对象
		2. 日志级别：Debug、Info、Warn、Error、DPanic、Panic、Fatal
		3. S函数和L函数很有用，提供了一个全局安全的logger对象，可以在任何地方使用
	*/
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger) //替换全局的logger对象
}
