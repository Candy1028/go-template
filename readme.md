```
your-project/
├── cmd/
│   └── server/
│       └── main.go         # 🚀 启动入口：加载配置、初始化服务、启动路由
├── config/
│   └── config.yaml          # ⚙️ 配置文件：数据库、MQ 连接参数、环境变量
├── internal/                # 🔒 应用核心逻辑（外部无法 import）
│   ├── middleware/          # 🛡️ HTTP 中间件（鉴权、日志、限流）
│   ├── model/               # 📦 数据结构定义（纯 struct）
│   │   └── user.go
│   ├── mqtopology/          # ✨ MQ 拓扑常量（队列名、交换机、路由键）
│   │   └── names.go
│   ├── consumer/            # 📨 MQ 消费者逻辑
│   │   ├── user.go          # 消费 user.registered 消息的业务处理
│   │   └── subscriber.go    # MQ 订阅、监听、重连封装
│   └── user/                # 👥 用户业务模块
│       ├── handler.go       # HTTP 接口层（参数校验、响应处理）
│       ├── repository.go    # 数据操作层（CRUD）
│       ├── service.go       # 核心业务逻辑
│       ├── routes.go        # 路由定义
│       └── publisher.go     # MQ 发布接口实现
├── pkg/                     # 📦 公共可复用库
│   ├── database/            # 🗃️ 数据库封装
│   │   └── mysql.go         # MySQL 连接池、查询封装
│   └── rabbitmq/            # 🔲 消息队列封装
│       └── client.go        # MQ 连接、Channel Pool、通用消费方法
│   └── util/                # 🧰 工具函数
│       └── jwt.go           # JWT 生成与验证
├── log/                     # 📜 日志配置与输出
├── test/                    # 🧪 测试相关
└── go.mod                   # 📌 Go Modules 依赖管理
```