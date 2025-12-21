//go:build !android && !linux
// +build !android,!linux

// Windows depends on Powershell and a bunch of stuff that is archived.  For that reason, we use the fallback position.
// See https://grahamwatts.co.uk/windows-secrets/

// MacOs, IoS. I don't have the tools to develop for those platforms.

package secrets

var ErrKeyringNotInstalled = ErrSecretsNotSupported
var ErrSecretManagerNotInstalled = ErrSecretsNotSupported

func NewSecretStore(appId, appDescription string) (Secretive, error) {
	return fallback()
}
