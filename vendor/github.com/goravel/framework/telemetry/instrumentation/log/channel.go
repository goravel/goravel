package log

import (
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/telemetry"
)

const defaultInstrumentationName = "github.com/goravel/framework/telemetry/instrumentation/log"

type TelemetryChannel struct{}

func NewTelemetryChannel() *TelemetryChannel {
	return &TelemetryChannel{}
}

func (r *TelemetryChannel) Handle(channelPath string) (log.Handler, error) {
	if telemetry.TelemetryFacade == nil {
		return nil, errors.TelemetryFacadeNotSet
	}

	config := telemetry.ConfigFacade
	if config == nil {
		return nil, errors.ConfigFacadeNotSet
	}

	instrumentName := config.GetString(channelPath+".instrument_name", defaultInstrumentationName)
	return &handler{
		logger: telemetry.TelemetryFacade.Logger(instrumentName),
	}, nil
}
