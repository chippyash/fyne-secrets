//go:build linux

package tests__test

import (
	"testing"

	"github.com/chippyash/fyne-secrets/secrets"
	"github.com/stretchr/testify/assert"
)

func TestLinux_NewSecretsStore(t *testing.T) {
	sut, err := secrets.NewSecretStore("app-id", "app description")
	assert.NoError(t, err)
	assert.Implements(t, (*secrets.Secretive)(nil), sut)
	assert.IsType(t, &secrets.GnomeSecretStore{}, sut)
}

func TestLinux_Functionality(t *testing.T) {
	sut, err := secrets.NewSecretStore("app-id", "app description")
	assert.NoError(t, err)
	err = sut.Store("test", []byte("test"))
	assert.NoError(t, err)
	b, err := sut.Exists("test")
	assert.NoError(t, err)
	assert.True(t, b)
	v, err := sut.Load("test")
	assert.NoError(t, err)
	assert.Equal(t, "test", string(v))
	err = sut.Delete("test")
	assert.NoError(t, err)
	b, err = sut.Exists("test")
	assert.NoError(t, err)
	assert.False(t, b)
}
