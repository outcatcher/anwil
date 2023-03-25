//go:build windows

package logging

import "log"

// GetDefaultLogger returns default stdout logger.
func GetDefaultLogger() *log.Logger {
	return log.Default()
}
