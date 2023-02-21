package handlers

import (
	"crypto/ed25519"
	"fmt"

	"github.com/outcatcher/anwil/config"
	"github.com/outcatcher/anwil/internal/auth"
)

type server struct {
	PrivateKey ed25519.PrivateKey
}

func newAPI(cfg *config.ServerConfiguration) (*server, error) {
	privateKey, err := auth.LoadPrivateKey(cfg.KeyStoragePath)
	if err != nil {
		return nil, fmt.Errorf("error creating new API instance: %w", err)
	}

	return &server{
		PrivateKey: privateKey,
	}, nil
}
