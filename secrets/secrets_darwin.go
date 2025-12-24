//go:build !ci && !ios && !wasm && !test_web_driver && !mobile && !noos && !tinygo

package secrets

import (
	"errors"

	"github.com/chippyash/go-keyring"
)

// NewSecretStore returns a Secretive implementation based on system prerequisites or falls back to a default option.
func NewSecretStore(appId, appDescription string) (Secretive, error) {
	return InitSecretStore(
		func(key string, value []byte) error {
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
