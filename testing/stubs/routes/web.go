package routes

import (
	"context"
	"encoding/json"
	nethttp "net/http"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	supportfile "github.com/goravel/framework/support/file"
)

func Web() {
	// ------------------
	// Test Route
	// ------------------
	facades.Route.Prefix("group1").Middleware(TestContextMiddleware()).Group(func(route1 route.Route) {
		route1.Prefix("group2").Middleware(TestContextMiddleware1()).Group(func(route2 route.Route) {
			route2.Get("/middleware/{id}", func(request http.Request) {
				facades.Response.Success().Json(http.Json{
					"id":   request.Input("id"),
					"ctx":  request.Context().Value("ctx").(string),
					"ctx1": request.Context().Value("ctx1").(string),
				})
			})
		})
		route1.Middleware(TestContextMiddleware2()).Get("/middleware/{id}", func(request http.Request) {
			facades.Response.Success().Json(http.Json{
				"id":   request.Input("id"),
				"ctx":  request.Context().Value("ctx").(string),
				"ctx1": request.Context().Value("ctx1").(string),
			})
		})
	})

	facades.Route.Get("/input/{id}", func(request http.Request) {
		facades.Response.Json(nethttp.StatusOK, http.Json{
			"id": request.Input("id"),
		})
	})

	facades.Route.Post("/input/{id}", func(request http.Request) {
		facades.Response.Success().Json(http.Json{
			"id": request.Input("id"),
		})
	})

	facades.Route.Put("/input/{id}", func(request http.Request) {
		facades.Response.Success().Json(http.Json{
			"id": request.Input("id"),
		})
	})

	facades.Route.Delete("/input/{id}", func(request http.Request) {
		facades.Response.Success().Json(http.Json{
			"id": request.Input("id"),
		})
	})

	facades.Route.Options("/input/{id}", func(request http.Request) {
		facades.Response.Success().Json(http.Json{
			"id": request.Input("id"),
		})
	})

	facades.Route.Patch("/input/{id}", func(request http.Request) {
		facades.Response.Success().Json(http.Json{
			"id": request.Input("id"),
		})
	})

	facades.Route.Any("/any/{id}", func(request http.Request) {
		facades.Response.Success().Json(http.Json{
			"id": request.Input("id"),
		})
	})

	facades.Route.Static("static", "./public")
	facades.Route.StaticFile("static-file", "./resources/logo.png")
	facades.Route.StaticFS("static-fs", nethttp.Dir("./public"))

	facades.Route.Middleware(TestAbortMiddleware()).Get("/middleware/{id}", func(request http.Request) {
		facades.Response.Success().Json(http.Json{
			"id": request.Input("id"),
		})
	})

	facades.Route.Middleware(TestContextMiddleware(), TestContextMiddleware1()).Get("/middlewares/{id}", func(request http.Request) {
		facades.Response.Success().Json(http.Json{
			"id":   request.Input("id"),
			"ctx":  request.Context().Value("ctx"),
			"ctx1": request.Context().Value("ctx1"),
		})
	})

	facades.Route.Prefix("prefix1").Prefix("prefix2").Get("input/{id}", func(request http.Request) {
		facades.Response.Success().Json(http.Json{
			"id": request.Input("id"),
		})
	})

	// ------------------
	// Test Request
	// ------------------
	facades.Route.Prefix("request").Group(func(route route.Route) {
		route.Get("/get/{id}", func(request http.Request) {
			facades.Response.Success().Json(http.Json{
				"id":       request.Input("id"),
				"name":     request.Query("name", "Hello"),
				"header":   request.Header("Hello", "World"),
				"method":   request.Method(),
				"path":     request.Path(),
				"url":      request.Url(),
				"full_url": request.FullUrl(),
				"ip":       request.Ip(),
			})
		})
		route.Get("/headers", func(request http.Request) {
			str, _ := json.Marshal(request.Headers())
			facades.Response.Success().String(string(str))
		})
		route.Post("/post", func(request http.Request) {
			facades.Response.Success().Json(http.Json{
				"name": request.Form("name", "Hello"),
			})
		})
		route.Post("/bind", func(request http.Request) {
			type Test struct {
				Name string
			}
			var test Test
			_ = request.Bind(&test)
			facades.Response.Success().Json(http.Json{
				"name": test.Name,
			})
		})
		route.Post("/file", func(request http.Request) {
			file, err := request.File("file")
			if err != nil {
				facades.Response.Success().String("get file error")
				return
			}
			path := "./resources/test.png"
			if err := file.Store(path); err != nil {
				facades.Response.Success().String("store file error: " + err.Error())
				return
			}

			facades.Response.Success().Json(http.Json{
				"exist": supportfile.Exist(path),
			})
		})
	})

	// ------------------
	// Test Response
	// ------------------
	facades.Route.Prefix("response").Group(func(route route.Route) {
		route.Get("/json", func(request http.Request) {
			facades.Response.Json(nethttp.StatusOK, http.Json{
				"id": "1",
			})
		})
		route.Get("/string", func(request http.Request) {
			facades.Response.String(nethttp.StatusCreated, "Goravel")
		})
		route.Get("/success/json", func(request http.Request) {
			facades.Response.Success().Json(http.Json{
				"id": "1",
			})
		})
		route.Get("/success/string", func(request http.Request) {
			facades.Response.Success().String("Goravel")
		})
		route.Get("/file", func(request http.Request) {
			facades.Response.File("./resources/logo.png")
		})
		route.Get("/download", func(request http.Request) {
			facades.Response.Download("./resources/logo.png", "1.png")
		})
		route.Get("/header", func(request http.Request) {
			facades.Response.Header("Hello", "goravel").String(nethttp.StatusOK, "Goravel")
		})
	})
}

func TestAbortMiddleware() http.Middleware {
	return func(request http.Request) {
		request.AbortWithStatus(nethttp.StatusNonAuthoritativeInfo)
		return
	}
}

func TestContextMiddleware() http.Middleware {
	return func(request http.Request) {
		ctx := request.Context()
		ctx = context.WithValue(ctx, "ctx", "Goravel")
		request.WithContext(ctx)

		request.Next()
	}
}

func TestContextMiddleware1() http.Middleware {
	return func(request http.Request) {
		ctx := request.Context()
		ctx = context.WithValue(ctx, "ctx1", "Hello")
		request.WithContext(ctx)

		request.Next()
	}
}

func TestContextMiddleware2() http.Middleware {
	return func(request http.Request) {
		ctx := request.Context()
		ctx = context.WithValue(ctx, "ctx2", "World")
		request.WithContext(ctx)

		request.Next()
	}
}
