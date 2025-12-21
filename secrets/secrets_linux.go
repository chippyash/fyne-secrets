//go:build linux

package secrets

// This will currently work for Gnome base distributions. See https://grahamwatts.co.uk/gnome-secrets/
// @todo Make it work for other distributions

// NewSecretStore returns a Secretive implementation based on system prerequisites or falls back to a default option.
func NewSecretStore(appId, appDescription string) (Secretive, error) {
	isGnome, gnomeErr := IsGnome()
	if isGnome && gnomeErr == nil {
		return NewGnomeSecretStore(appId, appDescription), nil
	}
	// @todo Check for other Linux variants
	isKeyctl, keyctlErr := IsKeyCtl()
	if isKeyctl && keyctlErr == nil {
		return NewKeyctlSecretStore(appId, appDescription), nil
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
