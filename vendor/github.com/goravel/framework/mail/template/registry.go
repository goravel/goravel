package template

import (
	"sync"

	"github.com/goravel/framework/contracts/config"
	contractsmail "github.com/goravel/framework/contracts/mail"
	"github.com/goravel/framework/errors"
)

const (
	DriverHTML   = "html"
	DriverCustom = "custom"
)

var engines sync.Map

// Get retrieves a cached mail template engine, creating it if it doesn't exist.
// This function is safe for concurrent use.
func Get(config config.Config) (contractsmail.Template, error) {
	engineName := config.GetString("mail.template.default", "html")

	if cached, ok := engines.Load(engineName); ok {
		return cached.(contractsmail.Template), nil
	}

	engine, err := createEngine(config, engineName)
	if err != nil {
		return nil, err
	}

	// Atomically load an existing engine or store the new one.
	// This prevents a race condition where two goroutines might create the same engine.
	// The one that gets to LoadOrStore first wins, and the other's created engine is discarded.
	actual, _ := engines.LoadOrStore(engineName, engine)

	return actual.(contractsmail.Template), nil
}

func createEngine(config config.Config, engineName string) (contractsmail.Template, error) {
	driver := config.GetString("mail.template.engines."+engineName+".driver", DriverHTML)
	switch driver {
	case DriverHTML:
		path := config.GetString("mail.template.engines."+engineName+".path", "resources/views/mail")
		return NewHtml(path), nil
	case DriverCustom:
		via := config.Get("mail.template.engines."+engineName+".via", "")
		if via == "" {
			return nil, errors.MailTemplateEngineViaRequired.Args(engineName)
		}

		switch v := via.(type) {
		case contractsmail.Template:
			return v, nil
		case func() (contractsmail.Template, error):
			engine, err := v()
			if err != nil {
				return nil, errors.MailTemplateEngineFactoryFailed.Args(engineName, err)
			}
			return engine, nil
		default:
			return nil, errors.MailTemplateEngineViaInvalid.Args(engineName)
		}
	default:
		return nil, errors.MailTemplateEngineNotSupported.Args(driver)
	}
}
