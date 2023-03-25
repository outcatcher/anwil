//go:build !windows

package logging

import (
	"log"
	"log/syslog"
)

// GetDefaultLogger returns syslogger with "anwil-migrate" prefix.
func GetDefaultLogger() *log.Logger {
	sysWriter, err := syslog.New(syslog.LOG_INFO, "anwil-api")
	if err != nil {
		return log.Default()
	}

	sysLogger := log.New(sysWriter, "", log.LstdFlags)

	return sysLogger
}
