/*
Package auth provides ways to authorize in API.
*/
package auth

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt/v4"
)

var (
	errUnexpectedSignMethod  = errors.New("unexpected signing method")
	errInvalidPrivateKeySize = errors.New("private key size is invalid")
)

// Claims - JWT payload contents.
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GeneratePrivateKey generates new ED25519 key.
func GeneratePrivateKey() (ed25519.PrivateKey, error) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("error generating key: %w", err)
	}

	return key, nil
}

// LoadPrivateKey loads ED25519 private key from path.
func LoadPrivateKey(privateKeyPath string) (ed25519.PrivateKey, error) {
	if privateKeyPath == "" {
		return GeneratePrivateKey() // generate key if none provided
	}

	keyFile, err := os.Open(filepath.Clean(privateKeyPath))
	if err != nil {
		return nil, fmt.Errorf("error opening private key file: %w", err)
	}

	defer func() {
		closeErr := keyFile.Close()
		if closeErr != nil {
			log.Println(closeErr)
		}
	}()

	keyData, err := io.ReadAll(keyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %w", err)
	}

	decodedData := make([]byte, hex.DecodedLen(len(keyData)))
	if _, err := hex.Decode(decodedData, keyData); err != nil {
		return nil, fmt.Errorf("error decoding loaded key file: %w", err)
	}

	return decodedData, nil
}

// ValidateToken validates token and return JWT payload data.
func ValidateToken(tokenString string, publicKey crypto.PublicKey) (*Claims, error) {
	claims := new(Claims)

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodEd25519)
		if !ok {
			return nil, fmt.Errorf("%w: %s", errUnexpectedSignMethod, token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error validating JWT: %w", err)
	}

	return claims, nil
}

// GenerateToken generates token with given claims.
func GenerateToken(claims *Claims, privateKey ed25519.PrivateKey) (string, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("error generating token: %w", errInvalidPrivateKeySize)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error creating signed string: %w", err)
	}

	return tokenString, nil
}
