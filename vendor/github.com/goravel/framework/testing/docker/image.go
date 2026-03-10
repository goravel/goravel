package docker

import (
	"context"
	"fmt"
	"time"

	contractsprocess "github.com/goravel/framework/contracts/process"
	contractsdocker "github.com/goravel/framework/contracts/testing/docker"
	"github.com/goravel/framework/errors"
	supportdocker "github.com/goravel/framework/support/docker"
	"github.com/goravel/framework/support/str"
)

type ImageDriver struct {
	config  contractsdocker.ImageConfig
	image   contractsdocker.Image
	process contractsprocess.Process
}

func NewImageDriver(image contractsdocker.Image, process contractsprocess.Process) *ImageDriver {
	return &ImageDriver{
		image:   image,
		process: process,
	}
}

func (r *ImageDriver) Build() error {
	if r.process == nil {
		return errors.ProcessFacadeNotSet.SetModule(errors.ModuleTesting)
	}

	command, exposedPorts := supportdocker.ImageToCommand(&r.image)
	res := r.process.Run(command)
	if res.Failed() {
		return errors.TestingImageBuildFailed.Args(r.image.Repository, res.Error())
	}

	containerID := str.Of(res.Output()).Squish().String()
	if containerID == "" {
		return errors.TestingImageNoContainerId.Args(r.image.Repository)
	}

	r.config = contractsdocker.ImageConfig{
		ContainerID:  containerID,
		ExposedPorts: exposedPorts,
	}

	return nil
}

func (r *ImageDriver) Config() contractsdocker.ImageConfig {
	return r.config
}

func (r *ImageDriver) Ready(fn func() error, durations ...time.Duration) error {
	duration := 1 * time.Minute
	if len(durations) > 0 {
		duration = durations[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return errors.TestingImageReadyTimeout.Args(r.image.Repository, duration)
		default:
			if err := fn(); err == nil {
				return nil
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func (r *ImageDriver) Shutdown() error {
	if r.config.ContainerID != "" {
		if res := r.process.Run(fmt.Sprintf("docker stop %s", r.config.ContainerID)); res.Failed() {
			return errors.TestingImageStopFailed.Args(r.image.Repository, res.Error())
		}
	}

	return nil
}
