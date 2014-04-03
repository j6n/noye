package logger

import (
	"os"

	"github.com/kdar/factorlog"
)

type logger struct {
	*factorlog.FactorLog
}

var log *logger

// Get returns a cache, or new logger
func Get() *logger {
	if log == nil {
		fmtr := factorlog.NewStdFormatter("\r%{Date} %{Time} -%{SEVERITY}- %{File}:%{Line} -- %{Message}")
		log = &logger{factorlog.New(os.Stdout, fmtr)}
	}

	return log
}
