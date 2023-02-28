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
facades.Route.Get("/", userController.Show)

// ORM
facades.Orm.Query().With("Author").First(&user)

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

- [ ] Hash
- [ ] Crypt
- [ ] Support Websocket
- [ ] Broadcasting
- [ ] Delay Queue
- [ ] Queue supports DB driver
- [ ] Notifications
- [ ] Optimize unit tests

## Documentation

Online documentation [https://www.goravel.dev](https://www.goravel.dev)

> To optimize the documentation, please submit a PR to the documentation
> repository [https://github.com/goravel/docs](https://github.com/goravel/docs)

## Contributors

This project exists thanks to all the people who contribute.

<a href="https://github.com/hwbrzzl" target="_blank"><img src="https://avatars.githubusercontent.com/u/24771476?v=4" width="48" height="48"></a>
<a href="https://github.com/merouanekhalili" target="_blank"><img src="https://avatars.githubusercontent.com/u/1122628?v=4" width="48" height="48"></a>
<a href="https://github.com/hongyukeji" target="_blank"><img src="https://avatars.githubusercontent.com/u/23145983?v=4" width="48" height="48"></a>
<a href="https://github.com/DevHaoZi" target="_blank"><img src="https://avatars.githubusercontent.com/u/115467771?v=4" width="48" height="48"></a>
<a href="https://github.com/sidshrivastav" target="_blank"><img src="https://avatars.githubusercontent.com/u/28773690?v=4" width="48" height="48"></a>

## Group

Welcome more discussion in Telegram.

[https://t.me/goravel](https://t.me/goravel)

<p align="left"><img src="https://www.goravel.dev/telegram.jpg" width="200"></p>

## License

The Goravel framework is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
