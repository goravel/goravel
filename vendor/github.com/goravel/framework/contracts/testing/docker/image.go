package docker

import "time"

type Image struct {
	Cmd          []string
	Env          []string
	ExposedPorts []string
	Repository   string
	Tag          string
	Args         []string
}

type ImageDriver interface {
	// Build the image.
	Build() error
	// Config gets the image configuration.
	Config() ImageConfig
	// Ready checks if the image is ready by the given function until the given duration, default is 1 minute.
	Ready(func() error, ...time.Duration) error
	// Shutdown the image.
	Shutdown() error
}

type ImageConfig struct {
	ContainerID  string
	ExposedPorts []string
}
