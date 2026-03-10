package convert

import (
	"github.com/goravel/framework/contracts/facades"
	"github.com/goravel/framework/support/path"
	"github.com/goravel/framework/support/str"
)

func BindingToFacade(binding string) string {
	for facade, b := range facades.FacadeToBinding {
		if b == binding {
			return facade
		}
	}

	return ""
}

func FacadeToBinding(facade string) string {
	return facades.FacadeToBinding[facade]
}

func FacadeToFilepath(facade string) string {
	return path.Facade(str.Of(facade).Snake().String() + ".go")
}
