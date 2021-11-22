package providers

import "goravel/routes"

type RouteServiceProvider struct {
}

func (router *RouteServiceProvider) Boot() {
	routes.V1()
}

func (router *RouteServiceProvider) Register() {

}
