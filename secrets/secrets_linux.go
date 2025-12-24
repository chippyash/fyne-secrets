//go:build linux

package secrets

import (
	"errors"

	"github.com/chippyash/go-keyring"
)

// NewSecretStore returns a Secretive implementation based on system prerequisites or falls back to a default option.
func NewSecretStore(appId, appDescription string) (Secretive, error) {
	isGnome, gnomeErr := IsGnome()
	if isGnome && gnomeErr == nil {
		return InitSecretStore(
			func(key string, value []byte) error {
				keyring.SetDescription(appDescription)
				return keyring.Set(appId, key, string(value))
			},
			func(key string) ([]byte, error) {
				s, err := keyring.Get(appId, key)
				return []byte(s), err
			},
			func(key string) (bool, error) {
				s, err := keyring.Get(appId, key)
				if err != nil {
					if errors.Is(err, keyring.ErrNotFound) {
						return false, nil
					}
					return false, err
				}
				if len(s) > 0 {
					return true, nil
				}
				return false, nil
			},
			func(key string) error {
				return keyring.Delete(appId, key)
			},
		), nil
	}

	// Check for other Linux variants
	isKeyctl, keyctlErr := IsKeyCtl()
	if isKeyctl && keyctlErr == nil {
		s, err := NewKeyctlSecretStore()
		if err == nil {
			return s, nil
		}
		keyctlErr = err
	}

	// there was an error, so we provide the fallback and the error
	if gnomeErr != nil {
		f, _ := fallback()
		return f, gnomeErr
	}
	if keyctlErr != nil {
		f, _ := fallback()
		return f, keyctlErr
	}

	//prerequisites not met
	return fallback()
}

// IsGnome checks if the GNOME Keyring is installed and returns true if it is, otherwise returns an error.
func IsGnome() (bool, error) {
	installed, err := PackageInstalled("gnome-keyring")
	if err != nil || !installed {
		return false, ErrKeyringNotInstalled
	}
	return true, nil
}
