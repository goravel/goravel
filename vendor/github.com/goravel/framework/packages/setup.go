package packages

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/packages"
	"github.com/goravel/framework/contracts/packages/modify"
	"github.com/goravel/framework/packages/options"
	"github.com/goravel/framework/packages/paths"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/env"
)

type setup struct {
	command     string
	driver      string
	facade      string
	force       bool
	onInstall   []modify.Apply
	onUninstall []modify.Apply
	paths       packages.Paths
}

var osExit = os.Exit

func Setup(args []string) packages.Setup {
	st := &setup{}
	var mainName string

	for _, arg := range args {
		if arg == "install" || arg == "uninstall" {
			st.command = arg
		}
		if arg == "--force" || arg == "-f" {
			st.force = true
		}
		if strings.HasPrefix(arg, "--facade=") {
			st.facade = strings.TrimPrefix(arg, "--facade=")
		}
		if strings.HasPrefix(arg, "--driver=") {
			st.driver = strings.TrimPrefix(arg, "--driver=")
		}
		if strings.HasPrefix(arg, "--main-path=") {
			mainName = strings.TrimPrefix(arg, "--main-path=")
		}
		if strings.HasPrefix(arg, "--paths=") {
			if err := json.Unmarshal([]byte(strings.TrimPrefix(arg, "--paths=")), &support.Config.Paths); err != nil {
				panic(fmt.Sprintf("failed to unmarshal paths: %s", err))
			}
		}
	}

	if mainName == "" {
		mainName = env.MainPath()
	}

	st.paths = Paths(mainName)

	return st
}

func (r *setup) Execute() {
	if r.command == "install" {
		for i := range r.onInstall {
			r.reportError(r.onInstall[i].Apply(options.Driver(r.driver), options.Force(r.force), options.Facade(r.facade)))
		}
	}

	if r.command == "uninstall" {
		for i := range r.onUninstall {
			r.reportError(r.onUninstall[i].Apply(options.Driver(r.driver), options.Force(r.force), options.Facade(r.facade)))
		}
	}
}

func (r *setup) Paths() packages.Paths {
	return r.paths
}

func (r *setup) Install(modifiers ...modify.Apply) packages.Setup {
	r.onInstall = modifiers

	return r
}

func (r *setup) Uninstall(modifiers ...modify.Apply) packages.Setup {
	r.onUninstall = modifiers

	return r
}

func (r *setup) reportError(err error) {
	if err != nil {
		if r.force {
			color.Warningln(err)
			return
		}

		color.Errorln(err)
		osExit(1)
	}
}

func Paths(mainPath ...string) packages.Paths {
	if len(mainPath) == 0 {
		mainPath = []string{env.MainPath()}
	}
	return paths.NewPaths(mainPath[0])
}
