<p align="center"><img src="https://user-images.githubusercontent.com/24771476/210227277-d2bbf608-1535-417a-98f0-a1103b813465.png" width="300"></p>

English | [中文](./README_zh.md)

## About Goravel

Goravel is a web application framework with complete functions and good scalability. As a starting scaffolding to help
Golang developers quickly build their own applications.

The framework style is consistent with [Laravel](https://github.com/laravel/laravel), let PHPer don't need to learn a
new framework, but also happy to play
around Golang! Tribute Laravel!

Welcome star, PR and issues！

## Getting started

```
// Generate APP_KEY
go run . artisan key:generate

// Route
facades.Route.Get("/", userController.Show)

// ORM
facades.Orm.Query().First(&user)

// Task Scheduling
facades.Schedule.Command("send:emails name").EveryMinute()

// Log
facades.Log.Debug(message)

// Cache
value := facades.Cache.Get("goravel", "default")

// Queues
err := facades.Queue.Job(&jobs.Test{}, []queue.Arg{}).Dispatch()
```

## Main Function

- [x] Config
- [x] Http
- [x] Authentication
- [x] Authorization
- [x] Orm
- [x] Migrate
- [x] Logger
- [x] Cache
- [x] Grpc
- [x] Artisan Console
- [x] Task Scheduling
- [x] Queue
- [x] Event
- [x] FileStorage
- [x] Mail
- [x] Validation
- [x] Mock

## Roadmap

- [ ] Optimize migration
- [ ] Orm relationships
- [ ] Custom .env path

## Documentation

Online documentation [https://www.goravel.dev](https://www.goravel.dev)

> To optimize the documentation, please submit a PR to the documentation
> repository [https://github.com/goravel/docs](https://github.com/goravel/docs)

## Group

Welcome more exchanges in Discord.

[https://discord.gg/cFc5csczzS](https://discord.gg/cFc5csczzS)

## License

The Goravel framework is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
