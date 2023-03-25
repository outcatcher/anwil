package config

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/outcatcher/anwil/domains/internals/config/schema"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
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

func generateYaml(config *schema.Configuration) string {
	return fmt.Sprintf(`
---
api:
  host: %s
  port: %d
`, config.API.Host, config.API.Port)
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

	ctx := context.Background()

	t.Run("Invalid config", func(t *testing.T) {
		t.Parallel()
		path := writeCfg(t, randomString("---", 120))

		_, err := LoadServerConfiguration(ctx, path)
		require.ErrorAs(t, err, new(*yaml.TypeError))
	})

	t.Run("Not existing file", func(t *testing.T) {
		t.Parallel()

		_, err := LoadServerConfiguration(ctx, randomString("path", 10))
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		expectedConfig := &schema.Configuration{
			API: schema.APIConfiguration{
				Host: randomString("host-", 20),
				Port: random.Intn(0xffff),
			},
			PrivateKeyPath: "",
		}

		fileName := writeCfg(t, generateYaml(expectedConfig))

		actual, err := LoadServerConfiguration(ctx, fileName)
		require.NoError(t, err)

		require.EqualValues(t, expectedConfig, actual)
	})
}
