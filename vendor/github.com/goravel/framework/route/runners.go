package route

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/route"
)

type RouteRunner struct {
	config config.Config
	route  route.Route
}

func NewRouteRunner(config config.Config, route route.Route) *RouteRunner {
	return &RouteRunner{
		config: config,
		route:  route,
	}
}

func (r *RouteRunner) Signature() string {
	return "route"
}

func (r *RouteRunner) ShouldRun() bool {
	return r.route != nil && r.config.GetString("http.default") != "" && r.config.GetBool("app.auto_run", true)
}

func (r *RouteRunner) Run() error {
	tlsHost := r.config.GetString("http.tls.host")
	tlsPort := r.config.GetString("http.tls.port")
	certFile := r.config.GetString("http.tls.ssl.cert")
	keyFile := r.config.GetString("http.tls.ssl.key")

	tlsShouldRun := tlsHost != "" && tlsPort != "" && certFile != "" && keyFile != ""
	if tlsShouldRun {
		if err := r.route.RunTLS(); err != nil {
			return err
		}
	}

	host := r.config.GetString("http.host")
	port := r.config.GetString("http.port")

	if host != "" && port != "" && (!tlsShouldRun || port != tlsPort) {
		if err := r.route.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (r *RouteRunner) Shutdown() error {
	return r.route.Shutdown()
}
