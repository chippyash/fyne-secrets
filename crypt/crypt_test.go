package crypt_test

import (
	"encoding/base64"
	"net/url"
	"testing"

	"github.com/chippyash/fyne-secrets/crypt"
	"github.com/stretchr/testify/assert"
)

func TestCrypt_GenKey(t *testing.T) {
	key, err := crypt.GenKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, key)
	u, err := url.QueryUnescape(key)
	assert.NoError(t, err)
	assert.NotEmpty(t, u)
	b64, err := base64.StdEncoding.DecodeString(u)
	assert.NoError(t, err)
	assert.NotEmpty(t, b64)
	assert.Equal(t, 32, len(b64))
}

func TestCrypt_NewWithoutKey(t *testing.T) {
	sut, err := crypt.NewCryptor(nil)
	assert.NoError(t, err)
	key := sut.Key()
	assert.NotEmpty(t, key)
}

func TestCrypt_NewWithKey(t *testing.T) {
	key, err := crypt.GenKey()
	assert.NoError(t, err)
	sut, err := crypt.NewCryptor(&key)
	assert.NoError(t, err)
	assert.Equal(t, key, sut.Key())
}

func TestCrypt_EncryptDecrypt(t *testing.T) {
	sut, err := crypt.NewCryptor(nil)
	assert.NoError(t, err)
	value := "35007634-bd2d-4b66-a159-15397617f6d"
	encrypted, err := sut.Encrypt(value)
	assert.NoError(t, err)
	assert.NotEqual(t, value, encrypted)
	decrypted, err := sut.Decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, value, decrypted)
}
