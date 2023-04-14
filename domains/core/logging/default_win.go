//go:build windows

package logging

// GetDefaultLogWriter returns writer to stdout.
func GetDefaultLogWriter() io.Writer {
	return os.Stdout
}
