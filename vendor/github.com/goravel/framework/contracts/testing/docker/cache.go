package docker

type CacheDriver interface {
	// Build a cache container, it doesn't wait for the cache to be ready, the Ready method needs to be called if
	// you want to check the container status.
	Build() error
	// Config get cache configuration.
	Config() CacheConfig
	// Fresh the cache.
	Fresh() error
	// Image gets the cache image.
	Image(image Image)
	// Ready checks if the cache is ready, the Build method needs to be called first.
	Ready() error
	// Reuse the existing cache container.
	Reuse(containerID string, port int) error
	// Shutdown the cache.
	Shutdown() error
}

type CacheConfig struct {
	ContainerID string
	Database    string
	Host        string
	Password    string
	Username    string
	Port        int
}
