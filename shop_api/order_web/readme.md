## 项目结构
```
├─order_web
│  ├─api
│  ├─config
│  ├─forms
│  ├─global
│  │  └─response
│  ├─initialize
│  ├─middlewares
│  ├─models
│  ├─proto
│  ├─router
│  ├─utils
│  ├─validator
│  ├─viper_test
│  └─zap_test
│      └─zap_log_file
```

- `config` : 配置
- `forms`：表单验证
- `global`： 全局变量
- `initialize`：初始化
- `middlerware`：自定义中间件
- `proto`
- `router`：路由
- `utils`：公用函数
- `validator`：自定义表达验证
- `main.go`： 入口

**注意**：
- config_debug.yaml是开发环境配置文件，
- initialize中的config.go中是配置文件加载的流程
- 配置环境变量`SHOP_DEBUG`为`true`,读取`config_debug.yaml`，否则读取`config_product.yaml`
