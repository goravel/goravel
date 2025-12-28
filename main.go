package main

import (
	"goravel/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	app.Wait()
}
