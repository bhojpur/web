package logger

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/bhojpur/web/pkg/client/logger/colors"
)

var errInvalidLogLevel = errors.New("web cli logger: invalid log level")

const (
	levelDebug = iota
	levelError
	levelFatal
	levelCritical
	levelSuccess
	levelWarn
	levelInfo
	levelHint
)

var (
	sequenceNo uint64
	instance   *CliLogger
	once       sync.Once
)
var debugMode = os.Getenv("DEBUG_ENABLED") == "1"

var logLevel = levelInfo

// CliLogger logs logging records to the specified io.Writer
type CliLogger struct {
	mu     sync.Mutex
	output io.Writer
}

// LogRecord represents a log record and contains the timestamp when the record
// was created, an increasing id, level and the actual formatted log line.
type LogRecord struct {
	ID       string
	Level    string
	Message  string
	Filename string
	LineNo   int
}

var Log = GetCliLogger(os.Stdout)

var (
	logRecordTemplate      *template.Template
	debugLogRecordTemplate *template.Template
)

// GetCliLogger initializes the logger instance with a NewColorWriter output
// and returns a singleton
func GetCliLogger(w io.Writer) *CliLogger {
	once.Do(func() {
		var (
			err             error
			simpleLogFormat = `{{Now "2006/01/02 15:04:05"}} {{.Level}} ▶ {{.ID}} {{.Message}}{{EndLine}}`
			debugLogFormat  = `{{Now "2006/01/02 15:04:05"}} {{.Level}} ▶ {{.ID}} {{.Filename}}:{{.LineNo}} {{.Message}}{{EndLine}}`
		)

		// Initialize and parse logging templates
		funcs := template.FuncMap{
			"Now":     Now,
			"EndLine": EndLine,
		}
		logRecordTemplate, err = template.New("simpleLogFormat").Funcs(funcs).Parse(simpleLogFormat)
		if err != nil {
			panic(err)
		}
		debugLogRecordTemplate, err = template.New("debugLogFormat").Funcs(funcs).Parse(debugLogFormat)
		if err != nil {
			panic(err)
		}

		instance = &CliLogger{output: colors.NewColorWriter(w)}
	})
	return instance
}

// SetOutput sets the logger output destination
func (l *CliLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = colors.NewColorWriter(w)
}

// Now returns the current local time in the specified layout
func Now(layout string) string {
	return time.Now().Format(layout)
}

// EndLine returns the a newline escape character
func EndLine() string {
	return "\n"
}

func (l *CliLogger) getLevelTag(level int) string {
	switch level {
	case levelFatal:
		return "FATAL   "
	case levelSuccess:
		return "SUCCESS "
	case levelHint:
		return "HINT    "
	case levelDebug:
		return "DEBUG   "
	case levelInfo:
		return "INFO    "
	case levelWarn:
		return "WARN    "
	case levelError:
		return "ERROR   "
	case levelCritical:
		return "CRITICAL"
	default:
		panic(errInvalidLogLevel)
	}
}

func (l *CliLogger) getColorLevel(level int) string {
	switch level {
	case levelCritical:
		return colors.RedBold(l.getLevelTag(level))
	case levelFatal:
		return colors.RedBold(l.getLevelTag(level))
	case levelInfo:
		return colors.BlueBold(l.getLevelTag(level))
	case levelHint:
		return colors.CyanBold(l.getLevelTag(level))
	case levelDebug:
		return colors.YellowBold(l.getLevelTag(level))
	case levelError:
		return colors.RedBold(l.getLevelTag(level))
	case levelWarn:
		return colors.YellowBold(l.getLevelTag(level))
	case levelSuccess:
		return colors.GreenBold(l.getLevelTag(level))
	default:
		panic(errInvalidLogLevel)
	}
}

// mustLog logs the message according to the specified level and arguments.
// It panics in case of an error.
func (l *CliLogger) mustLog(level int, message string, args ...interface{}) {
	if level > logLevel {
		return
	}
	// Acquire the lock
	l.mu.Lock()
	defer l.mu.Unlock()

	// Create the logging record and pass into the output
	record := LogRecord{
		ID:      fmt.Sprintf("%04d", atomic.AddUint64(&sequenceNo, 1)),
		Level:   l.getColorLevel(level),
		Message: fmt.Sprintf(message, args...),
	}

	err := logRecordTemplate.Execute(l.output, record)
	if err != nil {
		panic(err)
	}
}

// mustLogDebug logs a debug message only if debug mode
// is enabled. i.e. DEBUG_ENABLED="1"
func (l *CliLogger) mustLogDebug(message string, file string, line int, args ...interface{}) {
	if !debugMode {
		return
	}

	// Change the output to Stderr
	l.SetOutput(os.Stderr)

	// Create the log record
	record := LogRecord{
		ID:       fmt.Sprintf("%04d", atomic.AddUint64(&sequenceNo, 1)),
		Level:    l.getColorLevel(levelDebug),
		Message:  fmt.Sprintf(message, args...),
		LineNo:   line,
		Filename: filepath.Base(file),
	}
	err := debugLogRecordTemplate.Execute(l.output, record)
	if err != nil {
		panic(err)
	}
}

// Debug outputs a debug log message
func (l *CliLogger) Debug(message string, file string, line int) {
	l.mustLogDebug(message, file, line)
}

// Debugf outputs a formatted debug log message
func (l *CliLogger) Debugf(message string, file string, line int, vars ...interface{}) {
	l.mustLogDebug(message, file, line, vars...)
}

// Info outputs an information log message
func (l *CliLogger) Info(message string) {
	l.mustLog(levelInfo, message)
}

// Infof outputs a formatted information log message
func (l *CliLogger) Infof(message string, vars ...interface{}) {
	l.mustLog(levelInfo, message, vars...)
}

// Warn outputs a warning log message
func (l *CliLogger) Warn(message string) {
	l.mustLog(levelWarn, message)
}

// Warnf outputs a formatted warning log message
func (l *CliLogger) Warnf(message string, vars ...interface{}) {
	l.mustLog(levelWarn, message, vars...)
}

// Error outputs an error log message
func (l *CliLogger) Error(message string) {
	l.mustLog(levelError, message)
}

// Errorf outputs a formatted error log message
func (l *CliLogger) Errorf(message string, vars ...interface{}) {
	l.mustLog(levelError, message, vars...)
}

// Fatal outputs a fatal log message and exists
func (l *CliLogger) Fatal(message string) {
	l.mustLog(levelFatal, message)
	os.Exit(255)
}

// Fatalf outputs a formatted log message and exists
func (l *CliLogger) Fatalf(message string, vars ...interface{}) {
	l.mustLog(levelFatal, message, vars...)
	os.Exit(255)
}

// Success outputs a success log message
func (l *CliLogger) Success(message string) {
	l.mustLog(levelSuccess, message)
}

// Successf outputs a formatted success log message
func (l *CliLogger) Successf(message string, vars ...interface{}) {
	l.mustLog(levelSuccess, message, vars...)
}

// Hint outputs a hint log message
func (l *CliLogger) Hint(message string) {
	l.mustLog(levelHint, message)
}

// Hintf outputs a formatted hint log message
func (l *CliLogger) Hintf(message string, vars ...interface{}) {
	l.mustLog(levelHint, message, vars...)
}

// Critical outputs a critical log message
func (l *CliLogger) Critical(message string) {
	l.mustLog(levelCritical, message)
}

// Criticalf outputs a formatted critical log message
func (l *CliLogger) Criticalf(message string, vars ...interface{}) {
	l.mustLog(levelCritical, message, vars...)
}
