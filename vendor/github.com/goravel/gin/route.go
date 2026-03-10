package gin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/goravel/framework/contracts/config"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/str"
	"github.com/spf13/cast"
)

// map[path]map[method]info
var routes = make(map[string]map[string]contractshttp.Info)

var globalRecoverCallback func(ctx contractshttp.Context, err any) = defaultRecoverCallback

type Route struct {
	route.Router
	config           config.Config
	driver           string
	globalMiddleware []contractshttp.Middleware
	instance         *gin.Engine
	server           *http.Server
	tlsServer        *http.Server
}

func NewRoute(config config.Config, parameters map[string]any) (*Route, error) {
	driver := cast.ToString(parameters["driver"])
	if driver == "" {
		return nil, errors.New("please set the driver")
	}

	timeout := time.Duration(config.GetInt("http.request_timeout", 3)) * time.Second
	globalMiddleware := []contractshttp.Middleware{Timeout(timeout), Cors(), Tls()}

	route := &Route{
		config:           config,
		driver:           cast.ToString(parameters["driver"]),
		globalMiddleware: globalMiddleware,
	}
	if err := route.init(globalMiddleware); err != nil {
		return nil, err
	}

	return route, nil
}

func (r *Route) Fallback(handler contractshttp.HandlerFunc) {
	r.instance.NoRoute(handlerToGinHandler(handler))
}

func (r *Route) GetGlobalMiddleware() []contractshttp.Middleware {
	return r.globalMiddleware
}

func (r *Route) GetRoutes() []contractshttp.Info {
	paths := []string{}
	for path := range routes {
		paths = append(paths, path)
	}

	sort.Strings(paths)
	methods := []string{contractshttp.MethodGet + "|" + contractshttp.MethodHead, contractshttp.MethodHead, contractshttp.MethodGet, contractshttp.MethodPost, contractshttp.MethodPut, contractshttp.MethodDelete, contractshttp.MethodPatch, contractshttp.MethodOptions, contractshttp.MethodAny, contractshttp.MethodResource, contractshttp.MethodStatic, contractshttp.MethodStaticFile, contractshttp.MethodStaticFS}

	var infos []contractshttp.Info
	for _, path := range paths {
		for _, method := range methods {
			if info, ok := routes[path][method]; ok {
				infos = append(infos, info)
			}
		}
	}

	return infos
}

func (r *Route) GlobalMiddleware(middleware ...contractshttp.Middleware) {
	r.globalMiddleware = append(r.globalMiddleware, middleware...)
	if err := r.init(r.globalMiddleware); err != nil {
		panic(err)
	}
}

func (r *Route) Listen(l net.Listener) error {
	r.outputRoutes()
	color.Green().Println("[HTTP] Listening on: " + str.Of(l.Addr().String()).Start("http://").String())

	r.server = &http.Server{
		Addr:           l.Addr().String(),
		Handler:        http.AllowQuerySemicolons(r.instance),
		MaxHeaderBytes: r.config.GetInt(fmt.Sprintf("http.drivers.%s.header_limit", r.driver), 4096) << 10,
	}

	if err := r.server.Serve(l); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (r *Route) ListenTLS(l net.Listener) error {
	return r.ListenTLSWithCert(l, r.config.GetString("http.tls.ssl.cert"), r.config.GetString("http.tls.ssl.key"))
}

func (r *Route) ListenTLSWithCert(l net.Listener, certFile, keyFile string) error {
	r.outputRoutes()
	color.Green().Println("[HTTPS] Listening on: " + str.Of(l.Addr().String()).Start("https://").String())

	r.tlsServer = &http.Server{
		Addr:           l.Addr().String(),
		Handler:        http.AllowQuerySemicolons(r.instance),
		MaxHeaderBytes: r.config.GetInt(fmt.Sprintf("http.drivers.%s.header_limit", r.driver), 4096) << 10,
	}

	if err := r.tlsServer.ServeTLS(l, certFile, keyFile); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (r *Route) Info(name string) contractshttp.Info {
	routes := r.GetRoutes()

	for _, route := range routes {
		if route.Name == name {
			return route
		}
	}

	return contractshttp.Info{}
}

func (r *Route) Recover(callback func(ctx contractshttp.Context, err any)) {
	globalRecoverCallback = callback
	if err := r.init(r.globalMiddleware); err != nil {
		panic(err)
	}
}

func (r *Route) Run(host ...string) error {
	if len(host) == 0 {
		defaultHost := r.config.GetString("http.host")
		defaultPort := r.config.GetString("http.port")
		if defaultPort == "" {
			return errors.New("port can't be empty")
		}
		completeHost := defaultHost + ":" + defaultPort
		host = append(host, completeHost)
	}

	r.outputRoutes()
	color.Green().Println("[HTTP] Listening on: " + str.Of(host[0]).Start("http://").String())

	r.server = &http.Server{
		Addr:           host[0],
		Handler:        http.AllowQuerySemicolons(r.instance),
		MaxHeaderBytes: r.config.GetInt(fmt.Sprintf("http.drivers.%s.header_limit", r.driver), 4096) << 10,
	}

	if err := r.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (r *Route) RunTLS(host ...string) error {
	if len(host) == 0 {
		defaultHost := r.config.GetString("http.tls.host")
		defaultPort := r.config.GetString("http.tls.port")
		if defaultPort == "" {
			return errors.New("port can't be empty")
		}
		completeHost := defaultHost + ":" + defaultPort
		host = append(host, completeHost)
	}

	certFile := r.config.GetString("http.tls.ssl.cert")
	keyFile := r.config.GetString("http.tls.ssl.key")

	return r.RunTLSWithCert(host[0], certFile, keyFile)
}

func (r *Route) RunTLSWithCert(host, certFile, keyFile string) error {
	if host == "" {
		return errors.New("host can't be empty")
	}
	if certFile == "" || keyFile == "" {
		return errors.New("certificate can't be empty")
	}

	r.outputRoutes()
	color.Green().Println("[HTTPS] Listening on: " + str.Of(host).Start("https://").String())

	r.tlsServer = &http.Server{
		Addr:           host,
		Handler:        http.AllowQuerySemicolons(r.instance),
		MaxHeaderBytes: r.config.GetInt(fmt.Sprintf("http.drivers.%s.header_limit", r.driver), 4096) << 10,
	}

	if err := r.tlsServer.ListenAndServeTLS(certFile, keyFile); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (r *Route) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.instance.ServeHTTP(writer, request)
}

func (r *Route) SetGlobalMiddleware(middlewares []contractshttp.Middleware) {
	r.globalMiddleware = middlewares
	if err := r.init(r.globalMiddleware); err != nil {
		panic(err)
	}
}

func (r *Route) Shutdown(ctx ...context.Context) error {
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	if r.server != nil {
		return r.server.Shutdown(c)
	}
	if r.tlsServer != nil {
		return r.tlsServer.Shutdown(c)
	}
	return nil
}

func (r *Route) Test(request *http.Request) (*http.Response, error) {
	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	return recorder.Result(), nil
}

func (r *Route) init(globalMiddleware []contractshttp.Middleware) error {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableBindValidation()
	engine := gin.New()
	engine.MaxMultipartMemory = int64(r.config.GetInt(fmt.Sprintf("http.drivers.%s.body_limit", r.driver), 4096)) << 10

	ginMiddleware := []gin.HandlerFunc{}
	if r.config.GetBool("app.debug") {
		ginMiddleware = append(ginMiddleware, logMiddleware())
	}

	recoverMiddleware := func(ctx contractshttp.Context) {
		defer func() {
			if err := recover(); err != nil {
				globalRecoverCallback(ctx, err)
			}
		}()
		ctx.Request().Next()
	}
	globalMiddleware = append([]contractshttp.Middleware{recoverMiddleware}, globalMiddleware...)
	engine.Use(append(ginMiddleware, middlewaresToGinHandlers(globalMiddleware)...)...)

	template := r.config.Get("http.drivers." + r.driver + ".template")
	switch t := template.(type) {
	case render.HTMLRender:
		engine.HTMLRender = t
	case func() (render.HTMLRender, error):
		htmlRender, err := t()
		if err != nil {
			return err
		}

		engine.HTMLRender = htmlRender
	}

	if engine.HTMLRender == nil {
		var err error
		engine.HTMLRender, err = DefaultTemplate()
		if err != nil {
			return err
		}
	}

	r.Router = NewGroup(
		r.config,
		engine.Group("/"),
		"",
		[]contractshttp.Middleware{},
		[]contractshttp.Middleware{ResponseMiddleware()},
	)
	r.instance = engine

	return nil
}

func (r *Route) outputRoutes() {
	if r.config.GetBool("app.debug") && support.RuntimeMode != support.RuntimeArtisan && support.RuntimeMode != support.RuntimeTest {
		if err := App.MakeArtisan().Call("route:list"); err != nil {
			color.Errorln(fmt.Errorf("print route list failed: %w", err))
		}
	}
}

func defaultRecoverCallback(ctx contractshttp.Context, err any) {
	LogFacade.WithContext(ctx).Request(ctx.Request()).Error(err)
	ctx.Request().Abort(http.StatusInternalServerError)
}
