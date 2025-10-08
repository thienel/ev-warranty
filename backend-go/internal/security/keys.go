package security

import (
	"crypto/rsa"
	"fmt"
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	keyOnce    sync.Once
	keyErr     error
)

func InitRSAKeys(publicKeyPath, privateKeyPath string) error {
	keyOnce.Do(func() {
		publicKeyData, err := os.ReadFile(publicKeyPath)
		if err != nil {
			keyErr = fmt.Errorf("failed to read public key: %w", err)
			return
		}
		privateKeyData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			keyErr = fmt.Errorf("failed to read private key: %w", err)
			return
		}

		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
		if err != nil {
			keyErr = fmt.Errorf("failed to parse public key: %w", err)
			return
		}
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
		if err != nil {
			keyErr = fmt.Errorf("failed to parse private key: %w", err)
			return
		}
	})

	return keyErr
}

func PublicKey() *rsa.PublicKey   { return publicKey }
func PrivateKey() *rsa.PrivateKey { return privateKey }
