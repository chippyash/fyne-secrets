//go:build linux

package secrets

// This will currently work for Gnome base distributions. See https://grahamwatts.co.uk/gnome-secrets/
// @todo Make it work for non Gnome based distributions

// NewSecretStore returns a Secretive implementation based on system prerequisites or falls back to a default option.
func NewSecretStore(appId, appDescription string) (Secretive, error) {
	isGnome, err := IsGnome()
	if isGnome {
		return NewGnomeSecretStore(appId, appDescription), nil
	}
	// @todo Check for other Linux variants

	if err != nil {
		// there was an error, so we provide the fallback and the error
		f, _ := fallback()
		return f, err
	}

	//prerequisites not met
	return fallback()
}
