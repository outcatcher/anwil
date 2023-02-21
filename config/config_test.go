package config

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var random = rand.New(rand.NewSource(time.Now().UnixMicro())) //nolint:gosec

func randomString(prefix string, size int) string {
	bytes := make([]byte, size)

	_, _ = random.Read(bytes)

	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}

	message.NewPrinter(language.English)

	return prefix + string(bytes)
}

func generateYaml(config *ServerConfiguration) string {
	return fmt.Sprintf(`
---
host: %s
port: %d
`, config.Host, config.Port)
}

func writeCfg(t *testing.T, data string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "*-config.yaml")
	require.NoError(t, err)

	_, err = io.WriteString(tmpFile, data)
	require.NoError(t, err)

	return tmpFile.Name()
}

func TestLoadServerConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("Invalid config", func(t *testing.T) {
		t.Parallel()
		path := writeCfg(t, randomString("---", 120))

		_, err := LoadServerConfiguration(path)
		require.ErrorContains(t, err, "config decode error")
	})

	t.Run("Not existing file", func(t *testing.T) {
		t.Parallel()

		_, err := LoadServerConfiguration(randomString("path", 10))
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		expectedConfig := &ServerConfiguration{
			Host:           randomString("host-", 20),
			Port:           random.Intn(0xffff),
			KeyStoragePath: "",
		}

		fileName := writeCfg(t, generateYaml(expectedConfig))

		actual, err := LoadServerConfiguration(fileName)
		require.NoError(t, err)

		require.EqualValues(t, expectedConfig, actual)
	})
}
