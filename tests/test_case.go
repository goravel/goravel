package tests

import (
	"github.com/goravel/framework/testing"

	"goravel/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
