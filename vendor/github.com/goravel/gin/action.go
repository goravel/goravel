package gin

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	contractsroute "github.com/goravel/framework/contracts/route"
)

type Action struct {
	method string
	path   string
}

func NewAction(method, path, handler string) contractsroute.Action {
	if _, ok := routes[path]; !ok {
		routes[path] = make(map[string]contractshttp.Info)
	}

	if method == contractshttp.MethodGet {
		method = contractshttp.MethodGet + "|" + contractshttp.MethodHead
	}

	routes[path][method] = contractshttp.Info{
		Handler: handler,
		Method:  method,
		Path:    path,
	}

	return &Action{
		method: method,
		path:   path,
	}
}

func (r *Action) Name(name string) contractsroute.Action {
	info := routes[r.path][r.method]
	info.Name = name
	routes[r.path][r.method] = info

	return r
}
