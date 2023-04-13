//go:build !windows

package logging

import (
	"io"
	"log/syslog"
	"os"
)

// GetDefaultLogWriter returns writer to syslog with "anwil-migrate" prefix if possible.
func GetDefaultLogWriter() io.Writer {
	sysWriter, err := syslog.New(syslog.LOG_INFO, "anwil-api")
	if err != nil {
		return os.Stdout
	}

	return sysWriter
}
