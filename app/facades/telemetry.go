package facades

import (
	"github.com/goravel/framework/contracts/telemetry"
)

func Telemetry() telemetry.Telemetry {
	return App().MakeTelemetry()
}
