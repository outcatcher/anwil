/*
Package envyaml provide methods for using configuration from yaml and ENV variables.
*/
package envyaml

import (
	"context"
	"fmt"
	"io"

	"github.com/outcatcher/anwil/domains/internals/config/schema"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v2"
)

// Decode loads yaml configuration with some values overwritten with ENV variables.
func Decode(ctx context.Context, reader io.Reader, cfg *schema.Configuration) error {
	decoder := yaml.NewDecoder(reader)

	// first, load yaml contents
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("error decoding yaml: %w", err)
	}

	envCfg := new(schema.Configuration)

	// second, load required env vars
	if err := envconfig.Process(ctx, envCfg); err != nil {
		return fmt.Errorf("error loading env config: %w", err)
	}

	// overload only certain fields

	if envCfg.DB.Password != "" {
		cfg.DB.Password = envCfg.DB.Password
	}

	return nil
}
