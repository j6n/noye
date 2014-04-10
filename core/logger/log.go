package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/kdar/factorlog"
)

type logger struct {
	*factorlog.FactorLog
}

// New returns a new logger
func New() *logger {
	return &logger{factorlog.New(os.Stdout, NewScriptFormatter())}
}

var log *logger

// Get returns a cached or new logger
func Get() *logger {
	if log == nil {
		// 4 because caller -> lambda -> factorlog -> print
		log = &logger{factorlog.New(os.Stdout, NewNoyeFormatter(4))}
	}

	return log
}

type ScriptFormatter struct {
	*NoyeFormatter
}

func NewScriptFormatter() *ScriptFormatter {
	s := &ScriptFormatter{NewNoyeFormatter(0)}
	s.lines = func(int) string { return "" }
	return s
}

// NoyeFormatter is a factorlog formatter that applies the correct
// function name, line number and filename
type NoyeFormatter struct {
	temp  []byte
	lines func(n int) string
	depth int
}

// NewNoyeFormatter returns a new NoyeFormatter
func NewNoyeFormatter(depth int) *NoyeFormatter {
	return &NoyeFormatter{make([]byte, 64), getLines, depth}
}

// ShouldRuntimeCaller tells factorlog to use the runtime caller
func (n *NoyeFormatter) ShouldRuntimeCaller() bool {
	return false
}

// Format the message and produce a byte slice
func (n *NoyeFormatter) Format(context factorlog.LogContext) []byte {
	buf := &bytes.Buffer{}
	now := time.Now()

	year, month, day := now.Date()
	hour, min, sec := now.Clock()

	factorlog.NDigits(&n.temp, 4, 0, year)
	n.temp[4] = '/'
	factorlog.TwoDigits(&n.temp, 5, int(month))
	n.temp[7] = '/'
	factorlog.TwoDigits(&n.temp, 8, day)
	n.temp[10] = ' '

	factorlog.TwoDigits(&n.temp, 11, hour)
	n.temp[13] = ':'
	factorlog.TwoDigits(&n.temp, 14, min)
	n.temp[16] = ':'
	factorlog.TwoDigits(&n.temp, 17, sec)
	n.temp[19] = ' '

	buf.Write(n.temp[:20])

	buf.WriteString("[")
	buf.WriteString(factorlog.UcSeverityStrings[factorlog.SeverityToIndex(context.Severity)])
	buf.WriteString("] ")

	buf.WriteString(n.lines(n.depth))

	var message string
	if context.Format != nil {
		message = fmt.Sprintf(*context.Format, context.Args...)
	} else {
		message = fmt.Sprint(context.Args...)
	}

	buf.WriteString(message)
	length := len(message)
	if length > 0 && message[length-1] != '\n' {
		buf.WriteRune('\n')
	}

	return buf.Bytes()
}

func getLines(depth int) string {
	pc, file, line, _ := runtime.Caller(depth)
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	name := runtime.FuncForPC(pc).Name()
	name = name[strings.LastIndex(name, ".")+1:]
	return fmt.Sprintf("%s:%d#%s: ", short, line, name)
}
