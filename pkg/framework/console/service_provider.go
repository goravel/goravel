package console

type ServiceProvider struct {
}

func (console *ServiceProvider) Boot() {
	app := &Application{}
	app.Init()
}

func (console *ServiceProvider) Register() {
}
