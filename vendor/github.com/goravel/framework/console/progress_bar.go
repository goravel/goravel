package console

import (
	"github.com/pterm/pterm"

	"github.com/goravel/framework/contracts/console"
)

type ProgressBar struct {
	instance *pterm.ProgressbarPrinter
}

func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		instance: pterm.DefaultProgressbar.WithTotal(total).
			WithBarStyle(pterm.NewStyle(pterm.FgLightGreen)).
			WithTitleStyle(pterm.NewStyle(pterm.FgWhite)),
	}
}

func (r *ProgressBar) Advance(step ...int) {
	var instance *pterm.ProgressbarPrinter

	if len(step) > 0 {
		instance = r.instance.Add(step[0])
	} else {
		instance = r.instance.Increment()
	}

	if instance != nil {
		r.instance = instance
	}
}

func (r *ProgressBar) Finish() error {
	instance, err := r.instance.Stop()
	if err != nil {
		return err
	}
	r.instance = instance
	return nil
}

func (r *ProgressBar) SetTitle(message string) {
	r.instance = r.instance.UpdateTitle(message)
}

func (r *ProgressBar) ShowElapsedTime(b ...bool) console.Progress {
	r.instance = r.instance.WithShowElapsedTime(b...)
	return r
}

func (r *ProgressBar) ShowTitle(b ...bool) console.Progress {
	r.instance = r.instance.WithShowTitle(b...)
	return r
}

func (r *ProgressBar) Start() error {
	instance, err := r.instance.Start()
	if err != nil {
		return err
	}
	r.instance = instance
	return nil
}
