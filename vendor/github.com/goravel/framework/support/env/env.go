package env

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/goravel/framework/support"
)

// MainPath returns the package name of application, eg: goravel.
func MainPath() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Path
	}

	return "goravel"
}

// IsAir checks if the application is running using Air.
func IsAir() bool {
	for _, arg := range os.Args {
		if strings.Contains(filepath.ToSlash(arg), "/storage/temp") {
			return true
		}
	}

	return false
}

// IsArm returns whether the current CPU architecture is ARM.
// IsArm 返回当前 CPU 架构是否为 ARM。
func IsArm() bool {
	return runtime.GOARCH == "arm" || runtime.GOARCH == "arm64"
}

func IsArtisan() bool {
	return slices.Contains(os.Args, "artisan")
}

func IsBootstrapSetup() bool {
	data, err := os.ReadFile(filepath.Join(support.Config.Paths.Bootstrap, "app.go"))
	if err != nil {
		return false
	}

	return strings.Contains(string(data), "foundation.Setup().")
}

// IsDarwin returns whether the current operating system is Darwin.
// IsDarwin 返回当前操作系统是否为 Darwin。
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// IsDirectlyRun checks if the application is running using go run .
func IsDirectlyRun() bool {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Contains(filepath.Base(executable), os.TempDir()) ||
		(strings.Contains(filepath.ToSlash(executable), "/var/folders") && strings.Contains(filepath.ToSlash(executable), "/T/go-build")) // macOS
}

// IsGithub returns whether the current environment is github action.
// IsGithub 返回当前系统环境是否为 github action。
func IsGithub() bool {
	_, exists := os.LookupEnv("GITHUB_ACTION")

	return exists
}

// IsLinux returns whether the current operating system is Linux.
// IsLinux 返回当前操作系统是否为 Linux。
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsNoANSI checks if the application is running with the --no-ansi flag.
func IsNoANSI() bool {
	return slices.Contains(os.Args, "--no-ansi")
}

// IsTesting checks if the application is running in testing mode.
func IsTesting() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "-test.") {
			return true
		}
	}

	return false
}

// IsWindows returns whether the current operating system is Windows.
// IsWindows 返回当前操作系统是否为 Windows。
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsX86 returns whether the current CPU architecture is X86.
// IsX86 返回当前 CPU 架构是否为 X86。
func IsX86() bool {
	return runtime.GOARCH == "386" || runtime.GOARCH == "amd64"
}

// Is64Bit returns whether the current CPU architecture is 64-bit.
// Is64Bit 返回当前 CPU 架构是否为 64 位。
func Is64Bit() bool {
	return runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64"
}

func CurrentAbsolutePath() string {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, err := filepath.EvalSymlinks(filepath.Dir(executable))
	if err != nil {
		log.Fatal(err)
	}

	if IsTesting() || IsAir() || IsDirectlyRun() {
		res, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	return res
}
