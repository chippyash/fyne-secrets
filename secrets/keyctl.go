//go:build linux

package secrets

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

const (
	keyctlcmd = "keyctl"
)

type KeyctlSecretStore struct {
	cmd            string
	appIdentifier  string //The application Id - i.e. the one used to identify the Fyne app
	appDescription string //A description for the app. Used in the keychain
}

// NewKeyctlSecretStore initializes and returns a Secretive implementation using the keyctl Secret Tool as the backend.
func NewKeyctlSecretStore(appId, appDescription string) Secretive {
	//error can be ignored as it has already run successfully in the IsKeyCtl function
	cmd, _ := exec.LookPath(keyctlcmd)
	return &KeyctlSecretStore{
		cmd:            cmd,
		appIdentifier:  appId,
		appDescription: appDescription,
	}
}

// Store saves a secret value by associating it with a given key in the secret store using the `keyctl` command.
func (l *KeyctlSecretStore) Store(key string, value []byte) error {
	cmd := exec.Command(
		l.cmd,
		"add",
		"user",
		key,
		string(value),
		"@u",
	)
	cmdOutput, _ := cmd.StdoutPipe()
	_ = cmd.Start()
	cmdResult, err := io.ReadAll(cmdOutput)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading command output: %s", string(cmdResult)))
	}
	_ = cmd.Wait()
	return err
}

// Load retrieves the secret associated with the specified key from the secret store using the `keyctl` command.
func (l *KeyctlSecretStore) Load(key string) ([]byte, error) {
	cmd := exec.Command(
		l.cmd,
		"search",
		"@u",
		"user",
		key,
	)
	keyId, err := cmd.Output()
	if err != nil {
		return []byte{}, err
	}
	// the keyctl utility will return 1 if the secret doesn't exist
	if cmd.ProcessState.ExitCode() != 0 {
		return []byte{}, ErrSecretNotFound
	}

	cmd2 := exec.Command(
		l.cmd,
		"print",
		string(keyId),
	)
	return cmd2.Output()
}

// Exists checks whether a secret associated with the specified key exists in the secret store and returns true if found.
func (l *KeyctlSecretStore) Exists(key string) (bool, error) {
	cmd := exec.Command(
		l.cmd,
		"search",
		"@u",
		"user",
		key,
	)
	_, err := cmd.Output()
	if err != nil {
		return false, err
	}
	// the keyctl utility will return 1 if the secret doesn't exist
	if cmd.ProcessState.ExitCode() != 0 {
		return false, ErrSecretNotFound
	}

	return true, nil
}

// Delete removes the secret associated with the specified key from the secret store using the `keyctl` command.
func (l *KeyctlSecretStore) Delete(key string) error {
	cmd := exec.Command(
		l.cmd,
		"purge",
		"user",
		key,
	)
	_, err := cmd.Output()
	return err
}

// IsKeyCtl checks if the keyctl utility is installed and returns true if is, otherwise returns an error.
func IsKeyCtl() (bool, error) {
	_, err := PackageInstalled(keyctlcmd)
	if err != nil {
		return false, ErrSecretManagerNotInstalled
	}
	return true, nil
}
