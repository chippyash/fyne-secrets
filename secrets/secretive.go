package secrets

import (
	"fmt"

	"github.com/chippyash/fyne-secrets/secrets/fallbackstorage"
)

const (
	fileTpl = "data/%s.txt"
)

// Secretive is the interface to the secrets store.
// It is coincidentally the same as the fallbackstorage.Storing interface, but we keep it separate in case it diverges later.
type Secretive interface {
	// Store writes data
	Store(key string, value []byte) error
	// Load reads data
	Load(key string) ([]byte, error)
	// Exists checks if the data exists
	Exists(key string) (bool, error)
	// Delete deletes data
	Delete(key string) error
}

type StoreFunc func(key string, value []byte) error
type LoadFunc func(key string) ([]byte, error)
type ExistsFunc func(key string) (bool, error)
type DeleteFunc func(key string) error

type SecretStore struct {
	save   StoreFunc
	load   LoadFunc
	exists ExistsFunc
	delete DeleteFunc
}

// InitSecretStore returns a concrete implementation of Secretive
func InitSecretStore(s StoreFunc, l LoadFunc, e ExistsFunc, d DeleteFunc) Secretive {
	return &SecretStore{
		save:   s,
		load:   l,
		exists: e,
		delete: d,
	}
}

func (s *SecretStore) Store(key string, value []byte) error {
	return s.save(key, value)
}

func (s *SecretStore) Load(key string) ([]byte, error) {
	return s.load(key)
}

func (s *SecretStore) Exists(key string) (bool, error) {
	return s.exists(key)
}

func (s *SecretStore) Delete(key string) error {
	return s.delete(key)
}

// fallback initialises and returns a default Secretive implementation using file-based storage for secret management.
func fallback() (Secretive, error) {
	return InitSecretStore(
		func(key string, value []byte) error {
			fs := fallbackstorage.NewFileStorage()
			return fs.Store(fmt.Sprintf(fileTpl, key), value)
		},
		func(key string) ([]byte, error) {
			fs := fallbackstorage.NewFileStorage()
			return fs.Load(fmt.Sprintf(fileTpl, key))
		},
		func(key string) (bool, error) {
			fs := fallbackstorage.NewFileStorage()
			return fs.Exists(fmt.Sprintf(fileTpl, key))
		},
		func(key string) error {
			fs := fallbackstorage.NewFileStorage()
			return fs.Delete(fmt.Sprintf(fileTpl, key))
		},
	), nil
}
