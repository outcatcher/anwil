//go:build !windows

package logging

import (
	"log"
	"log/syslog"
	"os"
)

// syslog can be unavailable, e.g. in docker.
const syslogEnvVar = "SYSLOG"

// GetDefaultLogger returns syslogger with "anwil-migrate" prefix.
func GetDefaultLogger() *log.Logger {
	if os.Getenv(syslogEnvVar) == "" {
		return log.Default()
	}

	sysWriter, err := syslog.New(syslog.LOG_INFO, "anwil-api")
	if err != nil {
		log.Fatalln(err)
	}

	sysLogger := log.New(sysWriter, "", log.LstdFlags)

	return sysLogger
}
