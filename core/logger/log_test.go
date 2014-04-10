package logger

import (
	"os"
	"testing"

	"github.com/kdar/factorlog"
)

func TestLog(t *testing.T) {
	l := factorlog.New(os.Stdout, NewNoyeFormatter())
	l.Println("Custom formatter")
}
