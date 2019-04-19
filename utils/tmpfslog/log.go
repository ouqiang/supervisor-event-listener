package tmpfslog

import (
	"fmt"
	"os"
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

var f *os.File
var curLogLevel = INFO

func init() {
	fpath := "/tmp/supervisor-event-listener.log"
	f = newLogFile(fpath)
}

func newLogFile(fpath string) *os.File {
	mode := os.O_CREATE | os.O_WRONLY
	file, err := os.OpenFile(fpath, mode, 0644)
	if err != nil {
		panic(err)
	}
	return file
}

func log(level int, _fmt string, args ...interface{}) {
	if level > curLogLevel {
		return
	}
	now := time.Now()
	levelName := LEVELS_NAME[level]
	prefix := fmt.Sprintf("%s [%s]: ", now.Format(time.RFC3339), levelName)
	f.WriteString(prefix)
	f.WriteString(fmt.Sprintf(_fmt, args...))
	f.WriteString("\n")
}

func Debug(fmt string, args ...interface{}) {
	log(DEBUG, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	log(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	log(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	log(FATAL, fmt, args...)
	f.Close()
	os.Exit(127)
}
