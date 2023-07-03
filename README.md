<p align="center"><img src="https://www.goravel.dev/logo.png" width="300"></p>

English | [中文](./README_zh.md)

## About Goravel

Goravel is a web application framework with complete functions and good scalability. As a starting scaffolding to help
Gopher quickly build their own applications.

The framework style is consistent with [Laravel](https://github.com/laravel/laravel), let Phper don't need to learn a
new framework, but also happy to play around Golang! Tribute Laravel!

Welcome to star, PR and issues！

## Getting started

```
// Generate APP_KEY
go run . artisan key:generate

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

## Main Function

|             |                      |                      |                      |
| ----------  | --------------       | --------------       | --------------       |
| [Config](https://www.goravel.dev/getting-started/configuration.html)   | [Http](https://www.goravel.dev/the-basics/routing.html)  | [Authentication](https://www.goravel.dev/security/authentication.html)  | [Authorization](https://www.goravel.dev/security/authorization.html)  |
| [Orm](https://www.goravel.dev/ORM/getting-started.html)   | [Migrate](https://www.goravel.dev/ORM/migrations.html)  | [Logger](https://www.goravel.dev/the-basics/logging.html)  | [Cache](https://www.goravel.dev/digging-deeper/cache.html)  |
| [Grpc](https://www.goravel.dev/the-basics/grpc.html)   | [Artisan Console](https://www.goravel.dev/digging-deeper/artisan-console.html)  | [Task Scheduling](https://www.goravel.dev/digging-deeper/task-scheduling.html)  | [Queue](https://www.goravel.dev/digging-deeper/queues.html)  |
| [Event](https://www.goravel.dev/digging-deeper/event.html)   | [FileStorage](https://www.goravel.dev/digging-deeper/filesystem.html)  | [Mail](https://www.goravel.dev/digging-deeper/mail.html)  | [Validation](https://www.goravel.dev/the-basics/validation.html)  |
| [Mock](https://www.goravel.dev/digging-deeper/mock.html)   | [Hash](https://www.goravel.dev/security/hashing.html)  | [Crypt](https://www.goravel.dev/security/encryption.html)  | [Carbon](https://www.goravel.dev/digging-deeper/helpers.html)  |
| [Package Development](https://www.goravel.dev/digging-deeper/package-development.html)   |  |   |   |

## Roadmap

[For Detail](https://github.com/goravel/goravel/issues?q=is%3Aissue+is%3Aopen)

## Excellent Extend Packages

[For Detail](https://goravel.dev/prologue/packages.html)

## Contributors

This project exists thanks to all the people who contribute, to participate in the contribution, please see [Contribution Guide](https://goravel.dev/prologue/contributions.html).

<a href="https://github.com/hwbrzzl" target="_blank"><img src="https://avatars.githubusercontent.com/u/24771476?v=4" width="48" height="48"></a>
<a href="https://github.com/DevHaoZi" target="_blank"><img src="https://avatars.githubusercontent.com/u/115467771?v=4" width="48" height="48"></a>
<a href="https://github.com/merouanekhalili" target="_blank"><img src="https://avatars.githubusercontent.com/u/1122628?v=4" width="48" height="48"></a>
<a href="https://github.com/hongyukeji" target="_blank"><img src="https://avatars.githubusercontent.com/u/23145983?v=4" width="48" height="48"></a>
<a href="https://github.com/sidshrivastav" target="_blank"><img src="https://avatars.githubusercontent.com/u/28773690?v=4" width="48" height="48"></a>
<a href="https://github.com/Juneezee" target="_blank"><img src="https://avatars.githubusercontent.com/u/20135478?v=4" width="48" height="48"></a>
<a href="https://github.com/dragoonchang" target="_blank"><img src="https://avatars.githubusercontent.com/u/1432336?v=4" width="48" height="48"></a>
<a href="https://github.com/dhanusaputra" target="_blank"><img src="https://avatars.githubusercontent.com/u/35093673?v=4" width="48" height="48"></a>
<a href="https://github.com/mauri870" target="_blank"><img src="https://avatars.githubusercontent.com/u/10168637?v=4" width="48" height="48"></a>
<a href="https://github.com/Marian0" target="_blank"><img src="https://avatars.githubusercontent.com/u/624592?v=4" width="48" height="48"></a>
<a href="https://github.com/ahmed3mar" target="_blank"><img src="https://avatars.githubusercontent.com/u/12982325?v=4" width="48" height="48"></a>

## Group

Welcome more discussion in Telegram.

[https://t.me/goravel](https://t.me/goravel)

<p align="left"><img src="https://www.goravel.dev/telegram.jpg" width="200"></p>

## License

The Goravel framework is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
