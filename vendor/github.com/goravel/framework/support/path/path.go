package path

import (
	"github.com/goravel/framework/packages"
	packagespaths "github.com/goravel/framework/packages/paths"
	"github.com/goravel/framework/support"
)

func App(paths ...string) string {
	return packages.Paths().App().Abs(paths...)
}

func Base(paths ...string) string {
	return packagespaths.Abs(paths...)
}

func Bootstrap(paths ...string) string {
	return packages.Paths().Bootstrap().Abs(paths...)
}

func Config(paths ...string) string {
	return packages.Paths().Config().Abs(paths...)
}

func Database(paths ...string) string {
	return packages.Paths().Database().Abs(paths...)
}

func Executable(paths ...string) string {
	paths = append([]string{support.RootPath}, paths...)

	return Base(paths...)
}

func Facade(paths ...string) string {
	return packages.Paths().Facades().Abs(paths...)
}

func Lang(paths ...string) string {
	return packages.Paths().Lang().Abs(paths...)
}

func Migration(paths ...string) string {
	return packages.Paths().Migrations().Abs(paths...)
}

func Model(paths ...string) string {
	return packages.Paths().Models().Abs(paths...)
}

func Public(paths ...string) string {
	return packages.Paths().Public().Abs(paths...)
}

func Resource(paths ...string) string {
	return packages.Paths().Resources().Abs(paths...)
}

func Route(paths ...string) string {
	return packages.Paths().Routes().Abs(paths...)
}

func Storage(paths ...string) string {
	return packages.Paths().Storage().Abs(paths...)
}

func Test(paths ...string) string {
	return packages.Paths().Tests().Abs(paths...)
}

func View(paths ...string) string {
	return packages.Paths().Views().Abs(paths...)
}
