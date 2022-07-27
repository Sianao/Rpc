# Rpc 实现
# rpc

**C/S 模式**

- 服务注册

- 服务调用

- 心跳检测

- 负载均衡

- 超时检测

- 服务重连

- 服务扩展


## master

中心节点启动

`master.exe -p 端口号` 默认4563端口

## service

服务注册会将该服务所有方法 发送至**master** 节点登记 并与中心节点建立长连接

对于一个服务 可以发起多个链接 在服务运行的时候也可以动态添加服务实现无感拓展

```go
    s := api.NewService(":4563")
    // 注册服务 
    err := s.Register(&MS{}, 1)  //register 包含权重 默认为 1 用于负载均衡策略
    if err != nil {
        logrus.Error(err)
        return
    }
    s.Listen()
```

## client

客户端 与master建立tcp 连接 发送需要调用的服务以及方法和可选参数 等待服务段响应 一般情况存在超时响应 错误请求响应 以及正常响应 在发送请求的时候通过唯一ID进行标记 保证响应返回到正确的客户端

```go
    client, err := center.NewClient(":4563")
    if err != nil {
        panic(err)
    }
    re, err := client.Call("MS").FUN("Add")
    if err != nil {
        return
    }
    fmt.Println(re)
```