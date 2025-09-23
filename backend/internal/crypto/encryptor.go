package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

type Encryptor struct {
	gcm cipher.AEAD
	key []byte
}

type EncryptionConfig struct {
	Enabled   bool              `json:"enabled"`
	Algorithm string            `json:"algorithm"`
	KeyID     string            `json:"key_id"`
	Options   map[string]string `json:"options"`
}

func NewEncryptor() *Encryptor {
	return &Encryptor{}
}

func (e *Encryptor) Configure(config EncryptionConfig) error {
	if !config.Enabled {
		return nil
	}
	
	key := []byte("skylogix-storage-encryption-key-32")
	if len(key) != 32 {
		hash := sha256.Sum256([]byte(config.KeyID))
		key = hash[:]
	}
	
	e.key = key
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %w", err)
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}
	
	e.gcm = gcm
	return nil
}

func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	if e.gcm == nil {
		return data, nil // Encryption not configured
	}
	
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	ciphertext := e.gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (e *Encryptor) Decrypt(data []byte) ([]byte, error) {
	if e.gcm == nil {
		return data, nil // Encryption not configured
	}
	
	nonceSize := e.gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	
	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	
	return plaintext, nil
}

func (e *Encryptor) RotateKey(newKeyID string) error {
	
	config := EncryptionConfig{
		Enabled: true,
		KeyID:   newKeyID,
	}
	
	return e.Configure(config)
}

func (e *Encryptor) GetKeyInfo() map[string]interface{} {
	info := map[string]interface{}{
		"algorithm": "AES-256-GCM",
		"key_size":  len(e.key),
		"enabled":   e.gcm != nil,
	}
	
	return info
}
