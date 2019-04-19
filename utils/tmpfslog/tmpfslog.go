package tmpfslog

import (
	"fmt"
	"os"
	"time"
)

var f *os.File

func init() {
	fpath := "/tmp/supervisor-event-listener.log"
	// mode := os.O_CREATE | os.O_RDWR
	mode := os.O_APPEND | os.O_RDWR
	file, err := os.OpenFile(fpath, mode, 0644)
	if err != nil {
		panic(err)
	}
	f = file
}

func log(level int, _fmt string, args ...interface{}) {
	now := time.Now()
	f.WriteString(now.String())
	f.WriteString(fmt.Sprintf(_fmt, args...))
	f.WriteString("\n")
}

func Debug(fmt string, args ...interface{}) {
	log(3, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	log(3, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	log(3, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	log(3, fmt, args...)
}
