package data

import (
	"crypto/rsa"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	vault "github.com/hashicorp/vault/api"
)

// VaultConfig stores the Vault address and token.
type VaultConfig struct {
	Address string
	Token   string
}

// KeyManager is responsible for managing and rotating RSA keys loaded from Vault.
type KeyManager struct {
	vaultConfig  VaultConfig
	privateKey   *rsa.PrivateKey
	publicKeys   map[string]*rsa.PublicKey
	mu           sync.RWMutex
	refreshCycle time.Duration
	stopCh       chan struct{}
}

// NewKeyManager initializes a KeyManager and starts the key rotation mechanism.
func NewKeyManager(vaultConfig VaultConfig, refreshCycle time.Duration) (*KeyManager, error) {
	km := &KeyManager{
		vaultConfig:  vaultConfig,
		publicKeys:   make(map[string]*rsa.PublicKey),
		refreshCycle: refreshCycle,
		stopCh:       make(chan struct{}),
	}

	err := km.loadKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to load initial keys: %w", err)
	}

	go km.startKeyRotation()

	return km, nil
}

// loadKeys loads RSA private and public keys from Vault and stores them in KeyManager.
func (km *KeyManager) loadKeys() error {
	client, err := vault.NewClient(&vault.Config{
		Address: km.vaultConfig.Address,
	})
	if err != nil {
		return fmt.Errorf("failed to create Vault client: %w", err)
	}

	client.SetToken(km.vaultConfig.Token)

	// Load the private key
	secretPrivate, err := client.Logical().Read("jwt_keys/private_key")
	if err != nil {
		return fmt.Errorf("failed to read private key from Vault: %w", err)
	}
	if secretPrivate == nil {
		return fmt.Errorf("private key not found in Vault")
	}

	privateKeyData, ok := secretPrivate.Data["private_key"].(string)
	if !ok {
		return fmt.Errorf("unexpected format for private key data")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyData))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Load public keys (can be optimized with versioning, so we don't load all every time)
	secretPublic, err := client.Logical().Read("jwt_keys/public_keys")
	if err != nil {
		return fmt.Errorf("failed to read public keys from Vault: %w", err)
	}
	if secretPublic == nil {
		return fmt.Errorf("public keys not found in Vault")
	}

	publicKeyData, ok := secretPublic.Data["public_key"].(string)
	if !ok {
		return fmt.Errorf("unexpected format for public key data")
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyData))
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	km.mu.Lock()
	defer km.mu.Unlock()

	km.privateKey = privateKey
	km.publicKeys = map[string]*rsa.PublicKey{
		"default": publicKey, // Use a default kid for now
	}

	log.Println("Keys successfully loaded and updated from Vault")
	return nil
}

// startKeyRotation triggers the regular key refresh based on the defined refresh cycle.
func (km *KeyManager) startKeyRotation() {
	ticker := time.NewTicker(km.refreshCycle)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := km.loadKeys()
			if err != nil {
				log.Printf("Failed to refresh keys: %v. Retrying in 1 minute.", err)
				time.Sleep(time.Minute) // Retry after 1 minute
				continue
			}
		case <-km.stopCh:
			log.Println("Key rotation stopped")
			return
		}
	}
}

// Stop gracefully stops the key rotation process
func (km *KeyManager) Stop() {
	close(km.stopCh)
}

// GetPrivateKey safely retrieves the private key.
func (km *KeyManager) GetPrivateKey() *rsa.PrivateKey {
	km.mu.RLock()
	defer km.mu.RUnlock()

	return km.privateKey
}

// GetPublicKey safely retrieves a public key by its kid.
func (km *KeyManager) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	publicKey, exists := km.publicKeys[kid]
	if !exists {
		return nil, fmt.Errorf("public key not found for kid: %s", kid)
	}

	return publicKey, nil
}

// GetPublicKeys returns all public keys loaded in the manager.
func (km *KeyManager) GetPublicKeys() map[string]*rsa.PublicKey {
	km.mu.RLock()
	defer km.mu.RUnlock()

	return km.publicKeys
}
