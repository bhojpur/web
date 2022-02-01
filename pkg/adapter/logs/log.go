package logs

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

// Usage:
//
// import (
//	logs "github.com/bhojpur/logger/pkg/engine"
// )
//
//	log := NewLogger(10000)
//	log.SetLogger("console", "")
//
//	> the first params stand for how many channel
//
// Use it like this:
//
//	log.Trace("trace")
//	log.Info("info")
//	log.Warn("warning")
//	log.Debug("debug")
//	log.Critical("critical")

import (
	"log"
	"time"

	logs "github.com/bhojpur/logger/pkg/engine"
)

// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// levelLogLogger is defined to implement log.Logger
// the real log level will be LevelEmergency
const levelLoggerImpl = -1

// Name for adapter with Bhojpur official support
const (
	AdapterConsole   = "console"
	AdapterFile      = "file"
	AdapterMultiFile = "multifile"
	AdapterMail      = "smtp"
	AdapterConn      = "conn"
	AdapterEs        = "es"
	AdapterJianLiao  = "jianliao"
	AdapterSlack     = "slack"
	AdapterAliLS     = "alils"
)

// Legacy log level constants to ensure backwards compatibility.
const (
	LevelInfo  = LevelInformational
	LevelTrace = LevelDebug
	LevelWarn  = LevelWarning
)

type newLoggerFunc func() Logger

// Logger defines the behavior of a log provider.
type Logger interface {
	Init(config string) error
	WriteMsg(when time.Time, msg string, level int) error
	Destroy()
	Flush()
}

// Register makes a log provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, log newLoggerFunc) {
	logs.Register(name, func() logs.Logger {
		return &oldToNewAdapter{
			old: log(),
		}
	})
}

// BhojpurLogger is default logger in Bhojpur application.
// it can contain several providers and log message into all providers.
type BhojpurLogger logs.BhojpurLogger

const defaultAsyncMsgLen = 1e3

// NewLogger returns a new BhojpurLogger.
// channelLen means the number of messages in chan(used where asynchronous is true).
// if the buffering chan is full, logger adapters write to file or other way.
func NewLogger(channelLens ...int64) *BhojpurLogger {
	return (*BhojpurLogger)(logs.NewLogger(channelLens...))
}

// Async set the log to asynchronous and start the goroutine
func (bl *BhojpurLogger) Async(msgLen ...int64) *BhojpurLogger {
	(*logs.BhojpurLogger)(bl).Async(msgLen...)
	return bl
}

// SetLogger provides a given logger adapter into BhojpurLogger with config string.
// config need to be correct JSON as string: {"interval":360}.
func (bl *BhojpurLogger) SetLogger(adapterName string, configs ...string) error {
	return (*logs.BhojpurLogger)(bl).SetLogger(adapterName, configs...)
}

// DelLogger remove a logger adapter in BhojpurLogger.
func (bl *BhojpurLogger) DelLogger(adapterName string) error {
	return (*logs.BhojpurLogger)(bl).DelLogger(adapterName)
}

func (bl *BhojpurLogger) Write(p []byte) (n int, err error) {
	return (*logs.BhojpurLogger)(bl).Write(p)
}

// SetLevel Set log message level.
// If message level (such as LevelDebug) is higher than logger level (such as LevelWarning),
// log providers will not even be sent the message.
func (bl *BhojpurLogger) SetLevel(l int) {
	(*logs.BhojpurLogger)(bl).SetLevel(l)
}

// GetLevel Get Current log message level.
func (bl *BhojpurLogger) GetLevel() int {
	return (*logs.BhojpurLogger)(bl).GetLevel()
}

// SetLogFuncCallDepth set log funcCallDepth
func (bl *BhojpurLogger) SetLogFuncCallDepth(d int) {
	(*logs.BhojpurLogger)(bl).SetLogFuncCallDepth(d)
}

// GetLogFuncCallDepth return log funcCallDepth for wrapper
func (bl *BhojpurLogger) GetLogFuncCallDepth() int {
	return (*logs.BhojpurLogger)(bl).GetLogFuncCallDepth()
}

// EnableFuncCallDepth enable log funcCallDepth
func (bl *BhojpurLogger) EnableFuncCallDepth(b bool) {
	(*logs.BhojpurLogger)(bl).EnableFuncCallDepth(b)
}

// set prefix
func (bl *BhojpurLogger) SetPrefix(s string) {
	(*logs.BhojpurLogger)(bl).SetPrefix(s)
}

// Emergency Log EMERGENCY level message.
func (bl *BhojpurLogger) Emergency(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Emergency(format, v...)
}

// Alert Log ALERT level message.
func (bl *BhojpurLogger) Alert(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Alert(format, v...)
}

// Critical Log CRITICAL level message.
func (bl *BhojpurLogger) Critical(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Critical(format, v...)
}

// Error Log ERROR level message.
func (bl *BhojpurLogger) Error(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Error(format, v...)
}

// Warning Log WARNING level message.
func (bl *BhojpurLogger) Warning(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Warning(format, v...)
}

// Notice Log NOTICE level message.
func (bl *BhojpurLogger) Notice(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Notice(format, v...)
}

// Informational Log INFORMATIONAL level message.
func (bl *BhojpurLogger) Informational(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Informational(format, v...)
}

// Debug Log DEBUG level message.
func (bl *BhojpurLogger) Debug(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Debug(format, v...)
}

// Warn Log WARN level message.
// compatibility alias for Warning()
func (bl *BhojpurLogger) Warn(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Warn(format, v...)
}

// Info Log INFO level message.
// compatibility alias for Informational()
func (bl *BhojpurLogger) Info(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Info(format, v...)
}

// Trace Log TRACE level message.
// compatibility alias for Debug()
func (bl *BhojpurLogger) Trace(format string, v ...interface{}) {
	(*logs.BhojpurLogger)(bl).Trace(format, v...)
}

// Flush flush all chan data.
func (bl *BhojpurLogger) Flush() {
	(*logs.BhojpurLogger)(bl).Flush()
}

// Close close logger, flush all chan data and destroy all adapters in BhojpurLogger.
func (bl *BhojpurLogger) Close() {
	(*logs.BhojpurLogger)(bl).Close()
}

// Reset close all outputs, and set bl.outputs to nil
func (bl *BhojpurLogger) Reset() {
	(*logs.BhojpurLogger)(bl).Reset()
}

// GetBhojpurLogger returns the default BhojpurLogger
func GetBhojpurLogger() *BhojpurLogger {
	return (*BhojpurLogger)(logs.GetBhojpurLogger())
}

// GetLogger returns the default BhojpurLogger
func GetLogger(prefixes ...string) *log.Logger {
	return logs.GetLogger(prefixes...)
}

// Reset will remove all the adapter
func Reset() {
	logs.Reset()
}

// Async set the BhojpurLogger with Async mode and hold msglen messages
func Async(msgLen ...int64) *BhojpurLogger {
	return (*BhojpurLogger)(logs.Async(msgLen...))
}

// SetLevel sets the global log level used by the simple logger.
func SetLevel(l int) {
	logs.SetLevel(l)
}

// SetPrefix sets the prefix
func SetPrefix(s string) {
	logs.SetPrefix(s)
}

// EnableFuncCallDepth enable log funcCallDepth
func EnableFuncCallDepth(b bool) {
	logs.EnableFuncCallDepth(b)
}

// SetLogFuncCall set the CallDepth, default is 4
func SetLogFuncCall(b bool) {
	logs.SetLogFuncCall(b)
}

// SetLogFuncCallDepth set log funcCallDepth
func SetLogFuncCallDepth(d int) {
	logs.SetLogFuncCallDepth(d)
}

// SetLogger sets a new logger.
func SetLogger(adapter string, config ...string) error {
	return logs.SetLogger(adapter, config...)
}

// Emergency logs a message at emergency level.
func Emergency(f interface{}, v ...interface{}) {
	logs.Emergency(f, v...)
}

// Alert logs a message at alert level.
func Alert(f interface{}, v ...interface{}) {
	logs.Alert(f, v...)
}

// Critical logs a message at critical level.
func Critical(f interface{}, v ...interface{}) {
	logs.Critical(f, v...)
}

// Error logs a message at error level.
func Error(f interface{}, v ...interface{}) {
	logs.Error(f, v...)
}

// Warning logs a message at warning level.
func Warning(f interface{}, v ...interface{}) {
	logs.Warning(f, v...)
}

// Warn compatibility alias for Warning()
func Warn(f interface{}, v ...interface{}) {
	logs.Warn(f, v...)
}

// Notice logs a message at notice level.
func Notice(f interface{}, v ...interface{}) {
	logs.Notice(f, v...)
}

// Informational logs a message at info level.
func Informational(f interface{}, v ...interface{}) {
	logs.Informational(f, v...)
}

// Info compatibility alias for Warning()
func Info(f interface{}, v ...interface{}) {
	logs.Info(f, v...)
}

// Debug logs a message at debug level.
func Debug(f interface{}, v ...interface{}) {
	logs.Debug(f, v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
func Trace(f interface{}, v ...interface{}) {
	logs.Trace(f, v...)
}

func init() {
	SetLogFuncCallDepth(4)
}
