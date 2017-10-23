# Icecream-X


前言
---
**Icecream-X**是专为golang编写的高并发游戏服务端框架，**X**意为交错执行。


版本
---
version: 0.8


模块结构
---
- client (客户端模块，负责网络连接服务)
- server (服务端模块，负责服务端端口监听，与连接管理)
- proxy (负责代理终端的操作)
    - clientproxy (连接上服务端的客户端代理)
    - serverproxy (已连接上的服务端代理)
- parser (负责协议的解析并回调执行)
    - pbparser (用于处理protocbuf的协议模块)
    - gerpcparser (用于处理gevent rpc的协议模块)
- net_protocol (负责网络协议隔离)
    - tcp (用于tcp连接)
    - websocket (用于接入websocket)
