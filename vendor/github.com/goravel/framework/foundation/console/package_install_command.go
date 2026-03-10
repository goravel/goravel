package console

import (
	"fmt"
	"slices"
	"strings"

	contractsbinding "github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/facades"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/process"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/collect"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/convert"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/str"
)

type PackageInstallCommand struct {
	bindings      map[string]contractsbinding.Info
	chosenDrivers [][]contractsbinding.Driver
	paths         string
	process       process.Process
}

func NewPackageInstallCommand(bindings map[string]contractsbinding.Info, process process.Process, json contractsfoundation.Json) *PackageInstallCommand {
	paths, err := json.MarshalString(support.Config.Paths)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal paths: %s", err))
	}

	return &PackageInstallCommand{
		bindings: bindings,
		paths:    paths,
		process:  process,
	}
}

// Signature The name and signature of the console command.
func (r *PackageInstallCommand) Signature() string {
	return "package:install"
}

// Description The console command description.
func (r *PackageInstallCommand) Description() string {
	return "Install a package or a facade"
}

// Extend The console command extend.
func (r *PackageInstallCommand) Extend() command.Extend {
	return command.Extend{
		ArgsUsage: " <package@version> or <facade>",
		Category:  "package",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "all",
				Usage:   "Install all facades",
				Aliases: []string{"a"},
				Value:   false,
			},
			&command.BoolFlag{
				Name:    "default",
				Usage:   "Install facades with default drivers",
				Aliases: []string{"d"},
				Value:   false,
			},
			&command.BoolFlag{
				Name:  "dev",
				Usage: "Install drivers with the master branch",
				Value: false,
			},
		},
	}
}

// Handle Execute the console command.
func (r *PackageInstallCommand) Handle(ctx console.Context) error {
	names := ctx.Arguments()

	if len(names) == 0 {
		if ctx.OptionBool("all") {
			names = getAvailableFacades(r.bindings)
		} else {
			var err error

			options := []console.Choice{
				{Key: "All facades", Value: "all"},
				{Key: "Select facades", Value: "select"},
				{Key: "Third-party package", Value: "third"},
				{Key: "None", Value: "none"},
			}

			choice, err := ctx.Choice("Which facades or package do you want to install?", options)
			if err != nil {
				ctx.Error(err.Error())
				return nil
			}

			if choice == "all" {
				names = getAvailableFacades(r.bindings)
			}

			if choice == "select" {
				names, err = r.selectFacades(ctx)
			}

			if choice == "third" {
				var name string
				name, err = r.inputThirdPackage(ctx)
				if err == nil && name != "" {
					names = []string{name}
				}
			}

			if err != nil {
				ctx.Error(err.Error())
				return nil
			}
		}
	}

	for _, name := range names {
		if isPackage(name) {
			if err := r.installPackage(ctx, name); err != nil {
				ctx.Error(err.Error())
				return nil
			}
		} else {
			if doesFacadeExist(name) {
				if !ctx.OptionBool("all") {
					ctx.Info(fmt.Sprintf("Facade %s already exists", name))
				}
				continue
			}

			if err := r.installFacade(ctx, name); err != nil {
				ctx.Error(err.Error())
				return nil
			}
		}
	}

	// Clear chosen drivers after installation is complete, to avoid affecting subsequent installations.
	r.chosenDrivers = nil

	return nil
}

func (r *PackageInstallCommand) selectFacades(ctx console.Context) ([]string, error) {
	var facadeOptions []console.Choice
	for _, facade := range getAvailableFacades(r.bindings) {
		key := facade
		description := getFacadeDescription(facade, r.bindings)
		if description != "" {
			key = fmt.Sprintf("%-11s", facade) + color.Gray().Sprintf(" - %s", description)
		}
		facadeOptions = append(facadeOptions, console.Choice{
			Key:   key,
			Value: facade,
		})
	}

	return ctx.MultiSelect("Select the facades to install\nPlease check the how-to guide at the bottom", facadeOptions, console.MultiSelectOption{
		Filterable: true,
	})
}

func (r *PackageInstallCommand) inputThirdPackage(ctx console.Context) (string, error) {
	return ctx.Ask("Enter the package", console.AskOption{
		Description: "E.g.: github.com/goravel/framework or github.com/goravel/framework@master",
	})
}

func (r *PackageInstallCommand) installPackage(ctx console.Context, pkg string) error {
	if !strings.Contains(pkg, "@") && ctx.OptionBool("dev") {
		pkg += "@master"
	}

	pkgPath, _, _ := strings.Cut(pkg, "@")
	setup := pkgPath + "/setup"

	// get package
	if res := r.process.Run("go", "get", pkg); res.Failed() {
		return fmt.Errorf("failed to get package: %s", res.Error().Error())
	}

	// install package
	if res := r.process.WithSpinner("Installing "+pkg).Run("go", "run", setup, "install", "--main-path="+env.MainPath(), "--paths="+r.paths); res.Failed() {
		return fmt.Errorf("failed to install package: %s", res.Error().Error())
	}

	// tidy go.mod file
	if res := r.process.Run("go", "mod", "tidy"); res.Failed() {
		return fmt.Errorf("failed to tidy go.mod file: %s", res.Error().Error())
	}

	ctx.Success(fmt.Sprintf("Package %s installed successfully", pkg))

	return nil
}

func (r *PackageInstallCommand) installFacade(ctx console.Context, name string) error {
	binding := convert.FacadeToBinding(name)
	if _, exists := r.bindings[binding]; !exists {
		ctx.Warning(errors.PackageFacadeNotFound.Args(name).Error())
		ctx.Info(fmt.Sprintf("Available facades: %s", strings.Join(getAvailableFacades(r.bindings), ", ")))
		return nil
	}

	dependencyBindings := getDependencyBindings(binding, r.bindings, true)
	if len(dependencyBindings) > 0 && !ctx.OptionBool("all") {
		facades := make([]string, len(dependencyBindings))
		for i := range dependencyBindings {
			facades[i] = convert.BindingToFacade(dependencyBindings[i])
		}
		ctx.Info(fmt.Sprintf("%s depends on %s, they will be installed simultaneously", name, strings.Join(facades, ", ")))
	}

	bindingsToInstall := append(dependencyBindings, binding)
	for _, binding := range bindingsToInstall {
		facade := convert.BindingToFacade(binding)
		if doesFacadeExist(facade) {
			continue
		}

		bindingInfo := r.bindings[binding]
		setup := bindingInfo.PkgPath + "/setup"

		if res := r.process.WithSpinner("Installing "+facade).Run("go", "run", setup, "install", "--facade="+facade, "--main-path="+env.MainPath(), "--paths="+r.paths); res.Failed() {
			return fmt.Errorf("failed to install facade %s: %s", facade, res.Error().Error())
		}

		ctx.Success(fmt.Sprintf("Facade %s installed successfully", facade))

		if err := r.installDriver(ctx, facade, bindingInfo); err != nil {
			return err
		}

		if len(bindingInfo.Drivers) > 0 {
			r.chosenDrivers = append(r.chosenDrivers, bindingInfo.Drivers)
		}
	}

	if res := r.process.Run("go", "mod", "tidy"); res.Failed() {
		return fmt.Errorf("failed to tidy go.mod file: %s", res.Error().Error())
	}

	return nil
}

func (r *PackageInstallCommand) installDriver(ctx console.Context, facade string, bindingInfo contractsbinding.Info) error {
	if len(bindingInfo.Drivers) == 0 {
		return nil
	}

	// To avoid prompting for drivers that have already been installed
	for _, chooseDriver := range r.chosenDrivers {
		sortedChooseDriver := slices.Clone(chooseDriver)
		slices.SortFunc(sortedChooseDriver, func(a, b contractsbinding.Driver) int {
			return strings.Compare(a.Name, b.Name)
		})
		sortedDrivers := slices.Clone(bindingInfo.Drivers)
		slices.SortFunc(sortedDrivers, func(a, b contractsbinding.Driver) int {
			return strings.Compare(a.Name, b.Name)
		})
		if slices.Equal(sortedChooseDriver, sortedDrivers) {
			return nil
		}
	}

	var options []console.Choice
	for _, driver := range bindingInfo.Drivers {
		key := driver.Name
		if driver.Description != "" {
			key += color.Gray().Sprintf(" - %s", driver.Description)
		}

		options = append(options, console.Choice{
			Key:   key,
			Value: driver.Package,
		})
	}

	options = append(options, console.Choice{
		Key:   "Custom",
		Value: "Custom",
	})

	var driver string
	if ctx.OptionBool("default") {
		driver = options[0].Value
	} else {
		var err error
		driver, err = ctx.Choice(fmt.Sprintf("Select the %s driver to install", facade), options, console.ChoiceOption{
			Description: fmt.Sprintf("A driver is required for %s, please select one to install.", facade),
		})
		if err != nil {
			return err
		}

		if driver == "Custom" {
			driver, err = ctx.Ask(fmt.Sprintf("Please enter the %s driver package", facade))
			if err != nil {
				return err
			}
		}

		if driver == "" {
			return r.installDriver(ctx, facade, bindingInfo)
		}
	}

	if isInternalDriver(driver) {
		setup := bindingInfo.PkgPath + "/setup"

		if res := r.process.WithSpinner("Installing "+driver).Run("go", "run", setup, "install", "--driver="+driver, "--main-path="+env.MainPath(), "--paths="+r.paths); res.Failed() {
			return fmt.Errorf("failed to install driver %s: %s", driver, res.Error().Error())
		}

		ctx.Success(fmt.Sprintf("Driver %s installed successfully", driver))

		return nil
	}

	return r.installPackage(ctx, driver)
}

func doesFacadeExist(facade string) bool {
	return file.Exists(convert.FacadeToFilepath(facade))
}

func getAvailableFacades(bindings map[string]contractsbinding.Info) []string {
	var availableFacades []string
	for binding, info := range bindings {
		if !info.IsBase {
			availableFacades = append(availableFacades, convert.BindingToFacade(binding))
		}
	}

	slices.Sort(availableFacades)

	// Make sure "Route" facade is listed first, let the environment variables in .env.example be set up before other facades.
	targetIndex := -1
	for i, v := range availableFacades {
		if v == facades.Route {
			targetIndex = i
			break
		}
	}

	if targetIndex != -1 {
		value := availableFacades[targetIndex]
		availableFacades = append(availableFacades[:targetIndex], availableFacades[targetIndex+1:]...)
		availableFacades = append([]string{value}, availableFacades...)
	}

	return availableFacades
}

func getDependencyBindings(binding string, bindings map[string]contractsbinding.Info, withInstallTogether bool) []string {
	visited := make(map[string]bool)
	dependencyBindings := getDependencyBindingsRecursive(binding, bindings, visited, withInstallTogether)

	// binding may appear in InstallTogether, so it needs to be removed from dependencyBindings.
	dependencyBindings = collect.Filter(dependencyBindings, func(s string, _ int) bool {
		return s != binding
	})

	return dependencyBindings
}

func getDependencyBindingsRecursive(binding string, bindings map[string]contractsbinding.Info, visited map[string]bool, withInstallTogether bool) []string {
	var deps []string

	for _, dep := range bindings[binding].Dependencies {
		if info, ok := bindings[dep]; ok && !info.IsBase && !visited[dep] {
			visited[dep] = true
			deps = append(deps, getDependencyBindingsRecursive(dep, bindings, visited, withInstallTogether)...)
			deps = append(deps, dep)
		}
	}

	if withInstallTogether {
		for _, installTogetherBinding := range bindings[binding].InstallTogether {
			if info, ok := bindings[installTogetherBinding]; ok && !info.IsBase && !visited[installTogetherBinding] {
				visited[installTogetherBinding] = true
				deps = append(deps, getDependencyBindingsRecursive(installTogetherBinding, bindings, visited, withInstallTogether)...)
				deps = append(deps, installTogetherBinding)
			}
		}
	}

	return collect.Unique(deps)
}

func getFacadeDescription(facade string, bindings map[string]contractsbinding.Info) string {
	binding := convert.FacadeToBinding(facade)
	if info, exists := bindings[binding]; exists {
		return info.Description
	}

	return ""
}

func isPackage(pkg string) bool {
	return strings.Contains(pkg, "/")
}

func isInternalDriver(name string) bool {
	return name != "" && !str.Of(name).Contains(".", "/")
}
