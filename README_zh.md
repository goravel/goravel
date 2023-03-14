<p align="center"><img src="https://www.goravel.dev/logo.png" width="300"></p>

[English](./README.md) | 中文

# 关于 Goravel

Goravel 是一个功能完备、具有良好扩展能力的 Web 应用程序框架。作为一个起始脚手架帮助 Golang 开发者快速构建自己的应用。

框架风格与 [Laravel](https://github.com/laravel/laravel) 保持一致，让 Phper 不用学习新的框架，也可以愉快的玩转 Golang！致敬
Laravel！

欢迎 Star, PR, Issues！

## 快速上手

```
// 生成 APP_KEY
go run . artisan key:generate

// 定义路由
facades.Route.Get("/", userController.Show)

// 数据库查询
facades.Orm.Query().With("Author").First(&user)

// 任务调度
facades.Schedule.Command("send:emails name").EveryMinute()

// 记录 Log
facades.Log.Debug(message)

// 获取缓存
value := facades.Cache.Get("goravel", "default")

// 队列
err := facades.Queue.Job(&jobs.Test{}, []queue.Arg{}).Dispatch()
```

## 文档

在线文档 [https://www.goravel.dev/zh](https://www.goravel.dev/zh)

示例 [https://github.com/goravel/example](https://github.com/goravel/example)

> 优化文档，请提交 PR 至文档仓库 [https://github.com/goravel/docs](https://github.com/goravel/docs)

## 主要功能

- [x] 自定义配置
- [x] HTTP 服务
- [x] 用户认证
- [x] 用户授权
- [x] 数据库 ORM
- [x] 数据库迁移
- [x] 日志
- [x] 缓存
- [x] Grpc
- [x] Artisan 命令行
- [x] 任务调度
- [x] 队列
- [x] 事件系统
- [x] 文件存储
- [x] 邮件
- [x] 表单验证
- [x] Mock
- [x] Hash
- [x] Crypt

## 路线图

[查看详情](https://github.com/goravel/goravel/issues?q=is%3Aissue+is%3Aopen+label%3Aenhancement)

## 贡献者

这个项目的存在要归功于所有做出贡献的人。

<a href="https://github.com/hwbrzzl" target="_blank"><img src="https://avatars.githubusercontent.com/u/24771476?v=4" width="48" height="48"></a>
<a href="https://github.com/merouanekhalili" target="_blank"><img src="https://avatars.githubusercontent.com/u/1122628?v=4" width="48" height="48"></a>
<a href="https://github.com/hongyukeji" target="_blank"><img src="https://avatars.githubusercontent.com/u/23145983?v=4" width="48" height="48"></a>
<a href="https://github.com/DevHaoZi" target="_blank"><img src="https://avatars.githubusercontent.com/u/115467771?v=4" width="48" height="48"></a>
<a href="https://github.com/sidshrivastav" target="_blank"><img src="https://avatars.githubusercontent.com/u/28773690?v=4" width="48" height="48"></a>

## 群组

微信入群，请备注 Goravel

<p align="left"><img src="https://www.goravel.dev/wechat.jpg" width="200"></p>

## 开源许可

Goravel 框架是在 [MIT 许可](https://opensource.org/licenses/MIT) 下的开源软件。
