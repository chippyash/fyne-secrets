//go:build linux

package secrets

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

const (
	keyctlcmd = "keyctl"
)

type KeyctlSecretStore struct {
	cmd               string
	persistentKeyring string
}

// NewKeyctlSecretStore initializes and returns a Secretive implementation using the keyctl Secret Tool as the backend.
func NewKeyctlSecretStore() (Secretive, error) {
	//error can be ignored as it has already run successfully in the IsKeyCtl function
	tool, _ := exec.LookPath(keyctlcmd)
	cmd := exec.Command(
		tool,
		"get_persistent",
		"@u",
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	pkr := bytes.TrimRight(out, "\n")
	return &KeyctlSecretStore{
		cmd:               tool,
		persistentKeyring: string(pkr),
	}, nil
}

// Store saves a secret value by associating it with a given key in the secret store using the `keyctl` command.
func (l *KeyctlSecretStore) Store(key string, value []byte) error {
	cmd := exec.Command(
		l.cmd,
		"add",
		"user",
		key,
		string(value),
		l.persistentKeyring,
	)
	out, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading command output: %s", string(out)))
	}
	return nil
}

// Load retrieves the secret associated with the specified key from the secret store using the `keyctl` command.
func (l *KeyctlSecretStore) Load(key string) ([]byte, error) {
	cmd := exec.Command(
		l.cmd,
		"search",
		l.persistentKeyring,
		"user",
		key,
	)
	k, err := cmd.Output()
	if err != nil {
		return []byte{}, err
	}
	// the keyctl utility will return 1 if the secret doesn't exist
	if cmd.ProcessState.ExitCode() != 0 {
		return []byte{}, ErrSecretNotFound
	}

	keyId := strings.TrimRight(string(k), "\n")
	_ = keyId
	if err != nil {
		return []byte{}, err
	}
	cmd2 := exec.Command(
		l.cmd,
		"print",
		keyId,
	)
	out, err := cmd2.Output()
	out = bytes.TrimRight(out, "\n")
	return out, err
}

// Exists checks whether a secret associated with the specified key exists in the secret store and returns true if found.
func (l *KeyctlSecretStore) Exists(key string) (bool, error) {
	cmd := exec.Command(
		l.cmd,
		"search",
		l.persistentKeyring,
		"user",
		key,
	)
	_, err := cmd.Output()
	if err != nil {
		e := err.Error()
		if "exit status 1" == e {
			return false, nil
		}
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
	installed, err := PackageInstalled(keyctlcmd)
	if err != nil || !installed {
		return false, ErrSecretManagerNotInstalled
	}
	return true, nil
}
