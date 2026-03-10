package docker

import (
	"fmt"
	"math/rand"
	"net"
	"strings"

	"github.com/goravel/framework/contracts/testing/docker"
	"github.com/goravel/framework/errors"
)

func ExposedPort(exposedPorts []string, port string) string {
	for _, exposedPort := range exposedPorts {
		splitExposedPort := strings.Split(exposedPort, ":")
		if len(splitExposedPort) != 2 {
			continue
		}

		if splitExposedPort[1] != port && !strings.Contains(splitExposedPort[1], port+"/") {
			continue
		}

		return splitExposedPort[0]
	}

	return ""
}

func ImageToCommand(image *docker.Image) (command string, exposedPorts []string) {
	if image == nil {
		return "", nil
	}

	commands := []string{"docker", "run", "--rm", "-d"}
	if len(image.Env) > 0 {
		for _, env := range image.Env {
			commands = append(commands, "-e", env)
		}
	}

	var ports []string
	if len(image.ExposedPorts) > 0 {
		for _, port := range image.ExposedPorts {
			if !strings.Contains(port, ":") {
				port = fmt.Sprintf("%d:%s", ValidPort(), port)
			}
			ports = append(ports, port)
			commands = append(commands, "-p", port)
		}
	}

	commands = append(commands, fmt.Sprintf("%s:%s", image.Repository, image.Tag))

	if len(image.Args) > 0 {
		commands = append(commands, image.Args...)
	}

	if len(image.Cmd) > 0 {
		commands = append(commands, image.Cmd...)
	}

	return strings.Join(commands, " "), ports
}

// Used by TestContainer, to simulate the port is using.
var TestPortUsing = false

func IsPortUsing(port int) bool {
	if TestPortUsing {
		return true
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if l != nil {
		errors.Ignore(l.Close)
	}

	return err != nil
}

func ValidPort() int {
	for range 60 {
		random := rand.Intn(10000) + 10000
		if !IsPortUsing(random) {
			return random
		}
	}

	return 0
}
