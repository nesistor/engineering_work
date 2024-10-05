package keys

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

// KeyManager is responsible for managing and rotating RSA public keys loaded from Vault.
type KeyManager struct {
	vaultConfig     VaultConfig
	publicKeys      map[string]*rsa.PublicKey
	adminPublicKey  *rsa.PublicKey
	mu              sync.RWMutex
	refreshCycle    time.Duration
}

// NewKeyManager initializes a KeyManager and starts the key rotation mechanism.
func NewKeyManager(vaultConfig VaultConfig, refreshCycle time.Duration) (*KeyManager, error) {
	km := &KeyManager{
		vaultConfig:  vaultConfig,
		publicKeys:   make(map[string]*rsa.PublicKey),
		refreshCycle: refreshCycle,
	}

	err := km.loadKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to load initial keys: %w", err)
	}

	go km.startKeyRotation()

	return km, nil
}

// loadKeys loads RSA public keys from Vault and stores them in KeyManager.
func (km *KeyManager) loadKeys() error {
	client, err := vault.NewClient(&vault.Config{
		Address: km.vaultConfig.Address,
	})
	if err != nil {
		return fmt.Errorf("failed to create Vault client: %w", err)
	}

	client.SetToken(km.vaultConfig.Token)

	// Load public keys for users
	secretPublic, err := client.Logical().List("jwt_keys/public_keys")
	if err != nil {
		return fmt.Errorf("failed to list public keys from Vault: %w", err)
	}
	if secretPublic == nil {
		return fmt.Errorf("public keys not found in Vault")
	}

	publicKeys := make(map[string]*rsa.PublicKey)
	for kid, publicKeyData := range secretPublic.Data {
		keyDataStr, ok := publicKeyData.(string)
		if !ok {
			return fmt.Errorf("unexpected format for public key data for kid: %s", kid)
		}
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyDataStr))
		if err != nil {
			return fmt.Errorf("failed to parse public key for kid %s: %w", kid, err)
		}
		publicKeys[kid] = publicKey
	}

	// Load the admin public key
	secretAdminPublic, err := client.Logical().Read("jwt_keys/admin_public_key")
	if err != nil {
		return fmt.Errorf("failed to read admin public key from Vault: %w", err)
	}
	if secretAdminPublic == nil {
		return fmt.Errorf("admin public key not found in Vault")
	}

	adminPublicKeyData, ok := secretAdminPublic.Data["data"].(string)
	if !ok {
		return fmt.Errorf("unexpected format for admin public key data")
	}

	adminPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(adminPublicKeyData))
	if err != nil {
		return fmt.Errorf("failed to parse admin public key: %w", err)
	}

	km.mu.Lock()
	defer km.mu.Unlock()

	km.publicKeys = publicKeys
	km.adminPublicKey = adminPublicKey

	log.Println("Public keys successfully loaded and updated from Vault")
	return nil
}

// startKeyRotation triggers the regular key refresh based on the defined refresh cycle.
func (km *KeyManager) startKeyRotation() {
	ticker := time.NewTicker(km.refreshCycle)
	defer ticker.Stop()

	for range ticker.C {
		err := km.loadKeys()
		if err != nil {
			log.Printf("Failed to refresh keys: %v. Retrying in 1 minute.", err)
			time.Sleep(time.Minute) // Try to reload after one minute
			continue
		}
	}
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

// GetAdminPublicKey safely retrieves the admin public key.
func (km *KeyManager) GetAdminPublicKey() *rsa.PublicKey {
	km.mu.RLock()
	defer km.mu.RUnlock()

	return km.adminPublicKey
}

// GetPublicKeys returns all public keys loaded in the manager.
func (km *KeyManager) GetPublicKeys() map[string]*rsa.PublicKey {
	km.mu.RLock()
	defer km.mu.RUnlock()

	return km.publicKeys
}
