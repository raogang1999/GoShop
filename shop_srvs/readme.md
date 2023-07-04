

## user_srv

关于用户的微服务
```
 ─user_srv
    │  main.go
    │
    ├─global
    │      global.go
    │
    ├─handler
    │      user.go
    │
    ├─model
    │  │  user.go
    │  │
    │  └─main
    │          main.go
    │
    ├─proto
    │      user.pb.go
    │      user.proto
    │      user_grpc.pb.go
    │
    └─tests
            user.go
```

- main.go：启动微服务的入口
- global,全局对象创建与初始化
- handler，业务相关代码
- model，表结构
- proto，暴露的接口
    ```
    protoc --go_out=. --go-grpc_out=. -I="C:\Users\Rao\AppData\Local\JetBrains\GoLand2022.3\protoeditor" -I="proto" .\proto\user.proto   
    ```
- tests，测试代码