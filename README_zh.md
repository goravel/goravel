<p align="center"><img src="https://user-images.githubusercontent.com/24771476/210227277-d2bbf608-1535-417a-98f0-a1103b813465.png" width="300"></p>

[English](./README.md) | 中文

# 关于 Goravel

Goravel 是一个功能完备、具有良好扩展能力的 Web 应用程序框架。作为一个起始脚手架帮助 Golang 开发者快速构建自己的应用。

框架风格与 [Laravel](https://github.com/laravel/laravel) 保持一致，让 PHPer 不用学习新的框架，也可以愉快的玩转 Golang！致敬
Laravel！

欢迎 Star, PR, Issues！

## 快速上手

```
// 生成 APP_KEY
go run . artisan key:generate

// 定义路由
facades.Route.Get("/", userController.Show)

// 数据库查询
facades.Orm.Query().First(&user)

// 任务调度
facades.Schedule.Command("send:emails name").EveryMinute()

// 记录 Log
facades.Log.Debug(message)

// 获取缓存
value := facades.Cache.Get("goravel", "default")

// 队列
err := facades.Queue.Job(&jobs.Test{}, []queue.Arg{}).Dispatch()
```

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

## 路线图

- [ ] 优化迁移
- [ ] Orm 关联关系
- [ ] 自定义 .env 路径
- [ ] 数据库读写分离

## 文档

在线文档 [https://www.goravel.dev/zh](https://www.goravel.dev/zh)

> 优化文档，请提交 PR 至文档仓库 [https://github.com/goravel/docs](https://github.com/goravel/docs)

## 群组

欢迎在 Discord 中更多交流。

[https://discord.gg/cFc5csczzS](https://discord.gg/cFc5csczzS)

微信入群，请备注 Goravel

![](https://user-images.githubusercontent.com/24771476/194740900-cee4aa43-7c22-42b6-ada9-42bc160cd797.JPG)

## 开源许可

Goravel 框架是在 [MIT 许可](https://opensource.org/licenses/MIT) 下的开源软件。
