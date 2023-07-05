package main

import (
	"goravel/app/console"
	"goravel/bootstrap"
)

func main() {
	console.Newkernel(
		bootstrap.App(),
	).Handle()
}
