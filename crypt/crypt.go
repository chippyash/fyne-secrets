package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"net/url"
)

/*
 * Encryption/Decryption
 * @see https://dev.to/breda/secret-key-encryption-with-go-using-aes-316d for explanation of these crypt methods
 */

type Crypt interface {
	// Encrypt a string
	Encrypt(string) (string, error)
	// Decrypt a string
	Decrypt(string) (string, error)
	// Key returns the Key for encryption
	Key() string
}

type Cryptor struct {
	key *string
}

// Ensure that Cryptor implements the Crypt interface at compile time
var _ Crypt = &Cryptor{}

// NewCryptor creates a new Crypt implementation.
// If the key is nil, a new key will be generated. Retrieve the key with the Key() method
func NewCryptor(key *string) (Crypt, error) {
	if key == nil {
		k, err := GenKey()
		return &Cryptor{key: &k}, err
	}
	return &Cryptor{key: key}, nil
}

// Key returns the encryption key
func (c *Cryptor) Key() string {
	return *c.key
}

// Encrypt a string using AES-GCM
// Returns a url encoded base64 string
func (c *Cryptor) Encrypt(text string) (string, error) {
	key, err := decase(*c.key)
	if err != nil {
		return "", err
	}
	cphr, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)

	return encase(ciphertext), nil
}

// Decrypt a string using AES-GCM.
// `text` is a url encoded base64 string.
// Note: The nonce must be the first 12 bytes of the ciphertext.
func (c *Cryptor) Decrypt(text string) (string, error) {
	key, err := decase(*c.key)
	if err != nil {
		return "", err
	}
	cphr, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return "", err
	}

	cip, err := decase(text)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := cip[:nonceSize], cip[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// GenKey generates a new encryption key. 32-bit AES key encased as url encoded base64 string
func GenKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return encase(key), nil
}

// encase encodes a byte slice as url encoded base64
func encase(s []byte) string {
	return url.QueryEscape(base64.StdEncoding.EncodeToString(s))
}

// decase decodes a url encoded base64 string
func decase(s string) ([]byte, error) {
	u, err := url.QueryUnescape(s)
	if err != nil {
		return []byte(""), err
	}
	enc, err := base64.StdEncoding.DecodeString(u)
	if err != nil {
		return []byte(""), err
	}
	return enc, err
}
