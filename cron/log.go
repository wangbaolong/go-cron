package cron

import (
    "fmt"
    "time"
)

func Log(a ...interface{}) {
    fmt.Println(time.Now().Format("2006-01-02 15:04:05  "), a)
}

var logger Logger = new(defLogger)

// Logger is the interface used in this package for logging, so that any backend
// can be plugged in. It is a subset of the github.com/go-logr/logr interface.
type Logger interface {
    // Info logs routine messages about cron's operation.
    Info(msg string, keysAndValues ...interface{})
    // Error logs an error condition.
    Error(err error, msg string, keysAndValues ...interface{})
}

type defLogger struct {}

func (d *defLogger) Info(msg string, keysAndValues ...interface{}) {
    fmt.Println(time.Now().Format("2006-01-02 15:04:05 "), "Info", msg, keysAndValues)
}

func (d *defLogger) Error(err error, msg string, keysAndValues ...interface{}) {
    fmt.Println(time.Now().Format("2006-01-02 15:04:05 "), "Error", err, msg, keysAndValues)
}
