package secrets

import "github.com/pkg/errors"

var ErrSecretsNotSupported = errors.New("secrets not supported on this platform")
var ErrSecretNotFound = errors.New("secret not found")
var ErrKeyringNotInstalled = errors.New("keyring not installed, using fallback")
var ErrSecretManagerNotInstalled = errors.New("secrets manager not installed, using fallback")
