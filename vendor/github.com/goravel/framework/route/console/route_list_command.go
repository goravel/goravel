package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/support/collect"
)

var (
	routeParamsRegex = regexp.MustCompile(`({[^}]+})`)
	methodColors     = map[string]string{
		"ANY":         "<fg=red>ANY</>",
		"DELETE":      "<fg=red>DELETE</>",
		"GET":         "<fg=blue>GET</>",
		"HEAD":        "<fg=gray>HEAD</>",
		"OPTIONS":     "<fg=gray>OPTIONS</>",
		"POST":        "<fg=yellow>POST</>",
		"PUT":         "<fg=yellow>PUT</>",
		"PATCH":       "<fg=yellow>PATCH</>",
		"RESOURCE":    "<fg=green>RESOURCE</>",
		"STATIC":      "<fg=green>STATIC</>",
		"STATIC_FILE": "<fg=green>STATIC_FILE</>",
		"STATIC_FS":   "<fg=green>STATIC_FS</>",
	}
)

type RouteListCommand struct {
	router route.Route
}

func NewList(router route.Route) *RouteListCommand {
	return &RouteListCommand{
		router: router,
	}
}

// Signature The name and signature of the console command.
func (r *RouteListCommand) Signature() string {
	return "route:list"
}

// Description The console command description.
func (r *RouteListCommand) Description() string {
	return "List all registered routes"
}

// Extend The console command extend.
func (r *RouteListCommand) Extend() command.Extend {
	return command.Extend{
		Category: "route",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:  "method",
				Usage: "Filter the routes by method",
			},
			&command.StringFlag{
				Name:  "name",
				Usage: "Filter the routes by name",
			},
			&command.StringFlag{
				Name:  "path",
				Usage: "Filter the routes by path",
			},
			&command.StringSliceFlag{
				Name:  "except-path",
				Usage: "Do not display the routes matching the given path pattern",
			},
		},
	}
}

// Handle Execute the console command.
func (r *RouteListCommand) Handle(ctx console.Context) error {
	ctx.NewLine()
	routes := r.router.GetRoutes()
	if len(routes) == 0 {
		ctx.Warning("Your application doesn't have any routes.")
		return nil
	}

	filteredRoutes := filterRoutes(ctx, routes)
	if len(filteredRoutes) == 0 {
		ctx.Warning("Your application doesn't have any routes matching the given criteria.")
		return nil
	}

	for _, item := range filteredRoutes {
		ctx.TwoColumnDetail(fmt.Sprintf("%s %s", formatMethod(item.Method), formatPath(item.Path)), formateNameHandler(item.Name, item.Handler))
	}

	ctx.NewLine()
	ctx.TwoColumnDetail("", fmt.Sprintf("<fg=blue;op=bold>Showing [%d] routes</>", len(filteredRoutes)), ' ')

	return nil
}

func filterRoutes(ctx console.Context, routes []http.Info) []http.Info {
	var (
		matcher  []func(http.Info) bool
		contains = func(s, substr string) bool {
			return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
		}
	)
	if method := ctx.Option("method"); method != "" {
		matcher = append(matcher, func(route http.Info) bool {
			return contains(route.Method, method)
		})
	}

	if name := ctx.Option("name"); name != "" {
		matcher = append(matcher, func(route http.Info) bool {
			return contains(route.Name, name)
		})
	}

	if path := ctx.Option("path"); path != "" {
		matcher = append(matcher, func(route http.Info) bool {
			return contains(route.Path, path)
		})
	}

	if exceptPaths := ctx.OptionSlice("except-path"); len(exceptPaths) > 0 {
		matcher = append(matcher, func(route http.Info) bool {
			for _, exceptPath := range exceptPaths {
				if contains(route.Path, exceptPath) {
					return false
				}
			}

			return true
		})
	}

	return collect.Filter(routes, func(route http.Info, _ int) bool {
		for _, match := range matcher {
			if !match(route) {
				return false
			}
		}

		return true
	})
}

func formateNameHandler(name, handler string) string {
	if len(name) == 0 && len(handler) == 0 {
		return ""
	}

	if len(name) > 0 && len(handler) > 0 {
		name = name + " › "
	}

	return fmt.Sprintf("<fg=7472A3>%s%s</>", name, strings.TrimSuffix(handler, "-fm"))
}

func formatMethod(method string) string {
	split := strings.Split(method, "|")
	for i := range split {
		if colorized, ok := methodColors[split[i]]; ok {
			split[i] = colorized
		}
	}

	result := strings.Join(split, "<fg=gray>|</>")
	if padding := 12 - len(method); padding > 0 {
		result += strings.Repeat(" ", padding)
	}

	return result
}

func formatPath(path string) string {
	if cleared := strings.TrimPrefix(path, "/"); cleared != "" {
		path = cleared
	}

	return routeParamsRegex.ReplaceAllString(path, "<fg=yellow>$1</>")
}
