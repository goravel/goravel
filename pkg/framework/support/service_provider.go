package support

type ServiceProvider interface {
	Boot()
	Register()
}