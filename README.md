<div align="center">

<img src="https://www.goravel.dev/logo.png?v=1.14.x" width="300" alt="Logo">

[![Doc](https://pkg.go.dev/badge/github.com/goravel/framework)](https://pkg.go.dev/github.com/goravel/framework)
[![Go](https://img.shields.io/github/go-mod/go-version/goravel/framework)](https://go.dev/)
[![Release](https://img.shields.io/github/release/goravel/framework.svg)](https://github.com/goravel/framework/releases)
[![Test](https://github.com/goravel/framework/actions/workflows/test.yml/badge.svg)](https://github.com/goravel/framework/actions)
[![Report Card](https://goreportcard.com/badge/github.com/goravel/framework)](https://goreportcard.com/report/github.com/goravel/framework)
[![Codecov](https://codecov.io/gh/goravel/framework/branch/master/graph/badge.svg)](https://codecov.io/gh/goravel/framework)
![License](https://img.shields.io/github/license/goravel/framework)

</div>

English | [中文](./README_zh.md)

## About Goravel

Goravel is a web application framework with complete functions and good scalability. As a starting scaffolding to help
Gopher quickly build their own applications.

The framework style is consistent with [Laravel](https://github.com/laravel/laravel), let Php developer don't need to learn a
new framework, but also happy to play around Golang! In tribute to Laravel!

Welcome to star, PR and issues！

## Getting started

```
// Generate APP_KEY
./artisan key:generate

// Route
facades.Route().Get("/", userController.Show)

// ORM
facades.Orm().Query().With("Author").First(&user)

// Task Scheduling
facades.Schedule().Command("send:emails name").EveryMinute()

// Log
facades.Log().Debug(message)

// Cache
value := facades.Cache().Get("goravel", "default")

// Queues
err := facades.Queue().Job(&jobs.Test{}, []queue.Arg{}).Dispatch()
```

## Documentation

Online documentation [https://www.goravel.dev](https://www.goravel.dev)

Example [https://github.com/goravel/example](https://github.com/goravel/example)

> To optimize the documentation, please submit a PR to the documentation
> repository [https://github.com/goravel/docs](https://github.com/goravel/docs)

## Main Features

| Module Name | Description |
|-------------|-------------|
| [Artisan Console](https://www.goravel.dev/digging-deeper/artisan-console.html) | CLI command-line interface for application management and automation |
| [Authentication](https://www.goravel.dev/security/authentication.html) | User identity verification with JWT and Session drivers |
| [Authorization](https://www.goravel.dev/security/authorization.html) | Permission-based access control using policies and gates |
| [Cache](https://www.goravel.dev/digging-deeper/cache.html) | Store and retrieve data using memory, Redis, or custom drivers |
| [Carbon](https://www.goravel.dev/digging-deeper/helpers.html) | Helper functions for date and time manipulation |
| [Config](https://www.goravel.dev/getting-started/configuration.html) | Application configuration management from files and environment |
| [Crypt](https://www.goravel.dev/security/encryption.html) | Secure data encryption and decryption utilities |
| [DB](https://www.goravel.dev/database/getting-started.html) | Database query builder and connection management |
| [Event](https://www.goravel.dev/digging-deeper/event.html) | Application event dispatching and listening system |
| [Factory](https://www.goravel.dev/orm/factories.html) | Generate fake model data for testing purposes |
| [FileStorage](https://www.goravel.dev/digging-deeper/filesystem.html) | File upload, download, and storage across multiple drivers |
| [Grpc](https://www.goravel.dev/the-basics/grpc.html) | High-performance gRPC server and client implementation |
| [Hash](https://www.goravel.dev/security/hashing.html) | Secure password hashing using bcrypt algorithm |
| [Http](https://www.goravel.dev/the-basics/routing.html) | HTTP routing, controllers, and middleware management |
| [Http Client](https://www.goravel.dev/digging-deeper/http-client.html) | Make HTTP requests to external APIs and services |
| [Localization](https://www.goravel.dev/digging-deeper/localization.html) | Multi-language translation and locale management |
| [Logger](https://www.goravel.dev/the-basics/logging.html) | Application logging to files, console, or external services |
| [Mail](https://www.goravel.dev/digging-deeper/mail.html) | Send emails via SMTP or queue-based delivery |
| [Mock](https://www.goravel.dev/testing/mock.html) | Create test mocks for facades and dependencies |
| [Migrate](https://www.goravel.dev/database/migrations.html) | Version control for database schema changes |
| [Orm](https://www.goravel.dev/orm/getting-started.html) | Elegant Orm implementation for database operations |
| [Package Development](https://www.goravel.dev/digging-deeper/package-development.html) | Build reusable packages to extend framework functionality |
| [Queue](https://www.goravel.dev/digging-deeper/queues.html) | Defer time-consuming tasks to background job processing |
| [Seeder](https://www.goravel.dev/database/seeding.html) | Populate database tables with test or initial data |
| [Session](https://www.goravel.dev/the-basics/session.html) | Manage user session data across HTTP requests |
| [Task Scheduling](https://www.goravel.dev/digging-deeper/task-scheduling.html) | Schedule recurring tasks using cron-like expressions |
| [Testing](https://www.goravel.dev/testing/getting-started.html) | HTTP testing, mocking, and assertion utilities |
| [Validation](https://www.goravel.dev/the-basics/validation.html) | Validate incoming request data using rules |
| [View](https://www.goravel.dev/the-basics/views.html) | Template rendering engine for HTML responses |
| [TODO Process](https://www.goravel.dev/digging-deeper/process.html) | Long-running command-line process management |
| [TODO Telemetry](https://www.goravel.dev/digging-deeper/process.html) | Long-running command-line process management |

## Roadmap

[For Detail](https://github.com/goravel/goravel/issues?q=is%3Aissue+is%3Aopen)

## Excellent Extend Packages

[For Detail](https://www.goravel.dev/getting-started/packages.html)

## Contributors

This project exists thanks to all the people who contribute, to participate in the contribution, please see [Contribution Guide](https://www.goravel.dev/getting-started/contributions.html).

<a href="https://github.com/hwbrzzl" target="_blank"><img src="https://avatars.githubusercontent.com/u/24771476?v=4" width="48" height="48"></a>
<a href="https://github.com/DevHaoZi" target="_blank"><img src="https://avatars.githubusercontent.com/u/115467771?v=4" width="48" height="48"></a>
<a href="https://github.com/kkumar-gcc" target="_blank"><img src="https://avatars.githubusercontent.com/u/84431594?v=4" width="48" height="48"></a>
<a href="https://github.com/almas-x" target="_blank"><img src="https://avatars.githubusercontent.com/u/9382335?v=4" width="48" height="48"></a>
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
<a href="https://github.com/Kamandlou" target="_blank"><img src="https://avatars.githubusercontent.com/u/77993374?v=4" width="48" height="48"></a>
<a href="https://github.com/livghit" target="_blank"><img src="https://avatars.githubusercontent.com/u/108449432?v=4" width="48" height="48"></a>
<a href="https://github.com/jeff87218" target="_blank"><img src="https://avatars.githubusercontent.com/u/29706585?v=4" width="48" height="48"></a>
<a href="https://github.com/shayan-yousefi" target="_blank"><img src="https://avatars.githubusercontent.com/u/19957980?v=4" width="48" height="48"></a>
<a href="https://github.com/zxdstyle" target="_blank"><img src="https://avatars.githubusercontent.com/u/38398954?v=4" width="48" height="48"></a>
<a href="https://github.com/milwad-dev" target="_blank"><img src="https://avatars.githubusercontent.com/u/98118400?v=4" width="48" height="48"></a>
<a href="https://github.com/mdanialr" target="_blank"><img src="https://avatars.githubusercontent.com/u/48054961?v=4" width="48" height="48"></a>
<a href="https://github.com/KlassnayaAfrodita" target="_blank"><img src="https://avatars.githubusercontent.com/u/113383200?v=4" width="48" height="48"></a>
<a href="https://github.com/YlanzinhoY" target="_blank"><img src="https://avatars.githubusercontent.com/u/102574758?v=4" width="48" height="48"></a>
<a href="https://github.com/gouguoyin" target="_blank"><img src="https://avatars.githubusercontent.com/u/13517412?v=4" width="48" height="48"></a>
<a href="https://github.com/dzham" target="_blank"><img src="https://avatars.githubusercontent.com/u/10853451?v=4" width="48" height="48"></a>
<a href="https://github.com/praem90" target="_blank"><img src="https://avatars.githubusercontent.com/u/6235720?v=4" width="48" height="48"></a>
<a href="https://github.com/vendion" target="_blank"><img src="https://avatars.githubusercontent.com/u/145018?v=4" width="48" height="48"></a>
<a href="https://github.com/tzsk" target="_blank"><img src="https://avatars.githubusercontent.com/u/13273787?v=4" width="48" height="48"></a>
<a href="https://github.com/ycb1986" target="_blank"><img src="https://avatars.githubusercontent.com/u/12908032?v=4" width="48" height="48"></a>
<a href="https://github.com/BadJacky" target="_blank"><img src="https://avatars.githubusercontent.com/u/113529280?v=4" width="48" height="48"></a>
<a href="https://github.com/NiteshSingh17" target="_blank"><img src="https://avatars.githubusercontent.com/u/79739154?v=4" width="48" height="48"></a>
<a href="https://github.com/alfanzain" target="_blank"><img src="https://avatars.githubusercontent.com/u/4216529?v=4" width="48" height="48"></a>
<a href="https://github.com/oprudkyi" target="_blank"><img src="https://avatars.githubusercontent.com/u/3018472?v=4" width="48" height="48"></a>
<a href="https://github.com/zoryamba" target="_blank"><img src="https://avatars.githubusercontent.com/u/21248500?v=4" width="48" height="48"></a>
<a href="https://github.com/oguzhankrcb" target="_blank"><img src="https://avatars.githubusercontent.com/u/7572058?v=4" width="48" height="48"></a>
<a href="https://github.com/ChisThanh" target="_blank"><img src="https://avatars.githubusercontent.com/u/93512710?v=4" width="48" height="48"></a>
<a href="https://github.com/wyicwx" target="_blank"><img src="https://avatars.githubusercontent.com/u/1241187?v=4" width="48" height="48"></a>

## Sponsor

Better development of the project is inseparable from your support, reward us by [Open Collective](https://opencollective.com/goravel).

<p align="left"><img src="https://www.goravel.dev/reward.png" width="200"></p>

## Group

Welcome more discussion in Discord.

[https://discord.gg/cFc5csczzS](https://discord.gg/cFc5csczzS)

## License

The Goravel framework is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
