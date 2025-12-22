//go:build linux

package secrets

// NewSecretStore returns a Secretive implementation based on system prerequisites or falls back to a default option.
func NewSecretStore(appId, appDescription string) (Secretive, error) {
	isGnome, gnomeErr := IsGnome()
	if isGnome && gnomeErr == nil {
		return NewGnomeSecretStore(appId, appDescription), nil
	}
	// @todo Check for other Linux variants
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
