package errlog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	OFF = iota
	FATAL
	ERROR
	WARN
	INFO
	DEBUG
	ALL
)

var LEVELS_NAME = map[int]string{
	OFF:   "off",
	FATAL: "fatal",
	ERROR: "error",
	WARN:  "warn",
	INFO:  "info",
	DEBUG: "debug",
	ALL:   "all",
}

var f = os.Stderr
var curLogLevel = INFO

func init() {
	// fpath := "/tmp/supervisor-event-listener.log"
	// f = newLogFile(fpath)
}

func newLogFile(fpath string) *os.File {
	mode := os.O_CREATE | os.O_WRONLY
	file, err := os.OpenFile(fpath, mode, 0644)
	if err != nil {
		panic(err)
	}
	return file
}

func logFormat(level int, _fmt string, args ...interface{}) {
	if level > curLogLevel {
		return
	}

	_, fn, lineno, _ := runtime.Caller(2)
	fn = path.Base(fn)
	now := time.Now()
	levelName := LEVELS_NAME[level]
	prefix := fmt.Sprintf("%s [%s] %12s:%d ",
		now.Format(time.RFC3339), levelName, fn, lineno)
	f.WriteString(prefix)
	f.WriteString(fmt.Sprintf(_fmt, args...))
	f.WriteString("\n")
}

func Debug(fmt string, args ...interface{}) {
	logFormat(DEBUG, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	logFormat(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	logFormat(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	logFormat(FATAL, fmt, args...)
	f.Close()
	os.Exit(127)
}
