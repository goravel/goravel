<div align="center">

<img src="https://www.goravel.dev/logo.png" width="300" alt="Logo">

[![Doc](https://pkg.go.dev/badge/github.com/goravel/framework)](https://pkg.go.dev/github.com/goravel/framework)
[![Go](https://img.shields.io/github/go-mod/go-version/goravel/framework)](https://go.dev/)
[![Release](https://img.shields.io/github/release/goravel/framework.svg)](https://github.com/goravel/framework/releases)
[![Test](https://github.com/goravel/framework/actions/workflows/test.yml/badge.svg)](https://github.com/goravel/framework/actions)
[![Report Card](https://goreportcard.com/badge/github.com/goravel/framework)](https://goreportcard.com/report/github.com/goravel/framework)
[![Codecov](https://codecov.io/gh/goravel/framework/branch/master/graph/badge.svg)](https://codecov.io/gh/goravel/framework)
![License](https://img.shields.io/github/license/goravel/framework)

</div>

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
facades.Route().Get("/", userController.Show)

// 数据库查询
facades.Orm().Query().With("Author").First(&user)

// 任务调度
facades.Schedule().Command("send:emails name").EveryMinute()

// 记录 Log
facades.Log().Debug(message)

// 获取缓存
value := facades.Cache().Get("goravel", "default")

// 队列
err := facades.Queue().Job(&jobs.Test{}, []queue.Arg{}).Dispatch()
```

## 文档

在线文档 [https://www.goravel.dev/zh](https://www.goravel.dev/zh)

示例 [https://github.com/goravel/example](https://github.com/goravel/example)

> 优化文档，请提交 PR 至文档仓库 [https://github.com/goravel/docs](https://github.com/goravel/docs)

## 主要功能

|             |                      |                      |                      |
| ----------  | --------------       | --------------       | --------------       |
| [自定义配置](https://www.goravel.dev/zh/getting-started/configuration.html)   | [HTTP 服务](https://www.goravel.dev/zh/the-basics/routing.html)  | [用户认证](https://www.goravel.dev/zh/security/authentication.html)  | [用户授权](https://www.goravel.dev/zh/security/authorization.html)  |
| [数据库 ORM](https://www.goravel.dev/zh/ORM/getting-started.html)   | [数据库迁移](https://www.goravel.dev/zh/ORM/migrations.html)  | [日志](https://www.goravel.dev/zh/the-basics/logging.html)  | [缓存](https://www.goravel.dev/zh/digging-deeper/cache.html)  |
| [Grpc](https://www.goravel.dev/zh/the-basics/grpc.html)   | [Artisan 命令行](https://www.goravel.dev/zh/digging-deeper/artisan-console.html)  | [任务调度](https://www.goravel.dev/zh/digging-deeper/task-scheduling.html)  | [队列](https://www.goravel.dev/zh/digging-deeper/queues.html)  |
| [事件系统](https://www.goravel.dev/zh/digging-deeper/event.html)   | [文件存储](https://www.goravel.dev/zh/digging-deeper/filesystem.html)  | [邮件](https://www.goravel.dev/zh/digging-deeper/mail.html)  | [表单验证](https://www.goravel.dev/zh/the-basics/validation.html)  |
| [Mock](https://www.goravel.dev/zh/digging-deeper/mock.html)   | [Hash](https://www.goravel.dev/zh/security/hashing.html)  | [Crypt](https://www.goravel.dev/zh/security/encryption.html)  | [Carbon](https://www.goravel.dev/zh/digging-deeper/helpers.html)  |
| [扩展包开发](https://www.goravel.dev/zh/digging-deeper/package-development.html)   | [测试](https://www.goravel.dev/zh/testing/getting-started.html) |   |   |

## 路线图

[查看详情](https://github.com/goravel/goravel/issues?q=is%3Aissue+is%3Aopen)

## 优秀扩展包

[查看详情](https://goravel.dev/zh/prologue/packages.html)

## 贡献者

这个项目的存在要归功于所有做出贡献的人，参与贡献请查看[贡献指南](https://goravel.dev/zh/prologue/contributions.html)。

<a href="https://github.com/hwbrzzl" target="_blank"><img src="https://avatars.githubusercontent.com/u/24771476?v=4" width="48" height="48"></a>
<a href="https://github.com/DevHaoZi" target="_blank"><img src="https://avatars.githubusercontent.com/u/115467771?v=4" width="48" height="48"></a>
<a href="https://github.com/kkumar-gcc" target="_blank"><img src="https://avatars.githubusercontent.com/u/84431594?v=4" width="48" height="48"></a>
<a href="https://github.com/merouanekhalili" target="_blank"><img src="https://avatars.githubusercontent.com/u/1122628?v=4" width="48" height="48"></a>
<a href="https://github.com/hongyukeji" target="_blank"><img src="https://avatars.githubusercontent.com/u/23145983?v=4" width="48" height="48"></a>
<a href="https://github.com/sidshrivastav" target="_blank"><img src="https://avatars.githubusercontent.com/u/28773690?v=4" width="48" height="48"></a>
<a href="https://github.com/Juneezee" target="_blank"><img src="https://avatars.githubusercontent.com/u/20135478?v=4" width="48" height="48"></a>
<a href="https://github.com/dragoonchang" target="_blank"><img src="https://avatars.githubusercontent.com/u/1432336?v=4" width="48" height="48"></a>
<a href="https://github.com/dhanusaputra" target="_blank"><img src="https://avatars.githubusercontent.com/u/35093673?v=4" width="48" height="48"></a>
<a href="https://github.com/mauri870" target="_blank"><img src="https://avatars.githubusercontent.com/u/10168637?v=4" width="48" height="48"></a>
<a href="https://github.com/Marian0" target="_blank"><img src="https://avatars.githubusercontent.com/u/624592?v=4" width="48" height="48"></a>
<a href="https://github.com/ahmed3mar" target="_blank"><img src="https://avatars.githubusercontent.com/u/12982325?v=4" width="48" height="48"></a>
<a href="https://github.com/flc1125" target="_blank"><img src="https://avatars.githubusercontent.com/u/14297703?v=4" width="48" height="48"></a>
<a href="https://github.com/zzpwestlife" target="_blank"><img src="https://avatars.githubusercontent.com/u/12382180?v=4" width="48" height="48"></a>
<a href="https://github.com/juantarrel" target="_blank"><img src="https://avatars.githubusercontent.com/u/7213379?v=4" width="48" height="48"></a>

## 打赏

开源项目的发展离不开您的支持，感谢微信打赏。

<p align="left"><img src="https://www.goravel.dev/reward-wechat.jpg" width="200"></p>

## 群组

微信入群，请备注 Goravel

<p align="left"><img src="https://www.goravel.dev/wechat.jpg" width="200"></p>

## 开源许可

Goravel 框架是在 [MIT 许可](https://opensource.org/licenses/MIT) 下的开源软件。
