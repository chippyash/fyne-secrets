//go:build (!android && !linux) || generic
// +build !android,!linux generic

package tests__test

import (
	"testing"

	"github.com/chippyash/fyne-secrets/secrets"
	"github.com/stretchr/testify/assert"
)

func TestSecretive_InitSecretStore(t *testing.T) {
	sut := secrets.InitSecretStore(
		func(key string, value []byte) error {
			return nil
		},
		func(key string) ([]byte, error) {
			return []byte("test"), nil
		},
		func(key string) (bool, error) {
			return true, nil
		},
		func(key string) error {
			return nil
		},
	)

	assert.IsType(t, &secrets.SecretStore{}, sut)
	assert.Implements(t, (*secrets.Secretive)(nil), sut)
	assert.NoError(t, sut.Store("test", []byte("test")))
	v, err := sut.Load("test")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(v))
	b, err := sut.Exists("test")
	assert.NoError(t, err)
	assert.True(t, b)
	assert.NoError(t, sut.Delete("test"))
}

func TestSecrets_NewSecretsStore(t *testing.T) {
	sut, err := secrets.NewSecretStore("app-id", "app description")
	assert.NoError(t, err)
	assert.Implements(t, (*secrets.Secretive)(nil), sut)
}
