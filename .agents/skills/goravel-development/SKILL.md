---
name: goravel-development
description: >
  Use when writing or modifying Goravel application code: routing, controllers,
  facades, configuration, ORM, migrations, console commands, and tests. Covers
  both goravel/goravel (full) and goravel/goravel-lite scaffolds.
---

# Goravel Development

Laravel-style web framework for Go. Built on `github.com/goravel/framework`
plus a scaffold repo you clone and adapt. All repos live under
https://github.com/goravel (framework, scaffolds, example, installer, docs,
drivers).

Framework source: `github.com/goravel/framework`. Reference app:
`github.com/goravel/example`. Docs source: `github.com/goravel/docs`.

## Initialize & Run

```shell
go install github.com/goravel/installer/goravel@latest
goravel new blog
# Or clone manually:
#   Full scaffold:  git clone --depth=1 https://github.com/goravel/goravel.git
#   Lite scaffold:  git clone --depth=1 https://github.com/goravel/goravel-lite.git
cd <project> && go mod tidy && cp .env.example .env
./artisan key:generate                    # 32-char APP_KEY for encryption
./artisan jwt:secret                      # only if using Authentication
go run .                                  # start (air for live reload)
```

## Project Structure

The full scaffold installs every facade. The lite scaffold
ships a subset; install the rest with `./artisan package:install`. Add folders
freely, but don't rename defaults without `WithPaths()` in `bootstrap/app.go`.

- **`app/`** — Application code: HTTP controllers, console commands, models,
  middleware, events, listeners, jobs, AI agents/tools. `app/facades/` holds
  one thin file per facade delegating to `App().Make*()`.
- **`bootstrap/`** — App bootstrapping. `app.go` calls
  `foundation.Setup()...Create()`; `providers.go`, `commands.go`,
  `migrations.go` register providers, commands, migrations. `package:install`
  rewrites these as facades change.
- **`config/`** — One file per concern; each calls
  `facades.Config().Add("<name>", map[string]any{...})` in `init()`. Read env
  with `config.Env("KEY", default)`.
- **`database/`** — `migrations/` (schema, scaffold with `make:migration`),
  `seeders/` (initial/test data), `factories/` (fake model data for tests).
- **`routes/`** — Route definitions (`web.go`, `grpc.go`); registered in
  `bootstrap/app.go::WithRouting`.
- **`resources/views/`** — `*.tmpl` templates rendered via
  `ctx.Response().View().Make(...)`.
- **`lang/`**, **`storage/`**, **`public/`**, **`tests/`** — Translations,
  runtime files/logs, public assets, feature tests (`testify` suites).
- **`.env` / `artisan` / `main.go`** — Env config, console entry
  (`./artisan list`), app launch (`bootstrap.Boot().Start()`).

## Bootstrap & Configuration

`bootstrap/app.go` uses the `foundation.Setup()` builder; `With*` methods are
optional. `package:install` rewrites `bootstrap/app.go` and `providers.go`.

```go
return foundation.Setup().
    WithConfig(config.Boot).
    WithProviders(Providers).
    WithRouting(func() { routes.Web() }).
    WithCommands(Commands).
    WithMiddleware(func(h configuration.Middleware) { /* ... */ }).
    WithPaths(func(p configuration.Paths) { p.App("app") }).
    Create()
```

## Facades

Facades are the most important part of Goravel — all functions are implemented
via facades. Each facade resolves a container binding via `App().Make*()` at
call time, keeping access concise (`facades.Cache().Put(...)`) while staying
testable — swap the binding for a mock in tests (see Testing below).

All application facades live in `app/facades/`, one file per facade:

```go
func Cache() cache.Cache { return App().MakeCache() }
// usage: facades.Orm().Query().Get(&users)
```

### Install / Uninstall

The lite scaffold ships only `App`, `Artisan`, `Config`, `Process`. Add the rest
(rewrites `providers.go`, `config/`, `.env.example`, runs `go mod tidy`):

```shell
./artisan package:install Route --default    # one facade
./artisan package:install --all --default    # all facades + default drivers
./artisan package:uninstall Route
./artisan package:install github.com/goravel/redis  # external driver
```

> In the interactive picker, press `x` to select, then `Enter` to confirm.

### Finding Interfaces

Facade return types are in `github.com/goravel/framework/contracts/<module>/`. For GitHub,
read the version in `go.mod` (e.g. `v1.18.0`) and browse
`https://github.com/goravel/framework/tree/v1.18.x/contracts/<module>`
(minor version: `v1.18.0` → `v1.18.x`). Full facade list:
`github.com/goravel/framework/facades/facades.go`.

## Available Commands

```shell
./artisan list                    # all commands, args, options
./artisan make:<kind>             # scaffold (run list to see make:* generators)
```

From Go: `facades.Artisan().Call("list")`.

## Scaffolding Commands (make:*)

AI agents: always run `go run . artisan list | grep make:` to discover all
available scaffolding generators, then use them instead of hand-writing
boilerplate files. The available `make:*` commands depend on which facades
are installed — the list grows as you `package:install` more packages.

Run `go run . artisan make:<name> --help` to see flags and options for a
specific command before using it.

## Testing

`testify` suites + `tests.TestCase`. Scaffold: `./artisan make:test feature/UserTest`.

```go
type UserTestSuite struct {
    suite.Suite
    tests.TestCase
}
func TestUserTestSuite(t *testing.T) { suite.Run(t, new(UserTestSuite)) }
func (s *UserTestSuite) SetupTest() { s.RefreshDatabase() } // optional
func (s *UserTestSuite) TestIndex() { s.Seed(&seeders.UserSeeder{}) }
```

**Mocking** — use `github.com/goravel/framework/testing/mock`, never hand-roll:

```go
mockCache := mock.Factory().Cache()
mockCache.EXPECT().Get("name", "test").Return("Goravel").Once()
mockCache.AssertExpectations(t)
```

**Docker DB/Cache** — `facades.Testing().Docker().Database()` spins up an
isolated image per package for parallel `go test`. In `TestMain`:
`Build()` → `Ready()` → `Migrate()` → `m.Run()` → `Shutdown()`.

**Env files** — a `.env` in a test package dir overrides root `.env` for that
package. `go test ./... --env=.env.testing` swaps root env.

## Code Conventions

- `any` not `interface{}`. Standard Go formatting/naming; Go version in `go.mod`.
- Prefer `go test <pkg>` / `-run <Test>` over `go test ./...` (slow).
- `assert.*(t, *)` / `require.*(t, *)` directly, not `assert.New(t)`;
  `s.*(*, *)` in suites.
- Mocks: `.Once()` or `.Times(n)` always; `.Maybe()` only when necessary;
  avoid `mock.Anything` where a concrete value works.
- Table-driven tests with `t.Run` sub-tests.

## Where to Look

- **Reference app:** `github.com/goravel/example` — demonstrates every facade. Search first.
- **Docs:** https://www.goravel.dev (source `github.com/goravel/docs`).

## Custom Extensions

You can extend this skill with project-specific rules by creating
`.agents/skills/goravel-development/CUSTOM.md`. AI agents reading this
skill MUST also read and apply that file when it exists.
