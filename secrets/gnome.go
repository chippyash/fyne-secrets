//go:build linux

package secrets

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

const (
	gnomekeyringcmd = "gnome-keyring"
	gnomesecretcmd  = "secret-tool"
)

type GnomeSecretStore struct {
	cmd            string
	appIdentifier  string //The application Id - i.e. the one used to identify the Fyne app
	appDescription string //A description for the app. Used in the keychain
}

// NewGnomeSecretStore initializes and returns a Secretive implementation using the GNOME Secret Tool as the backend.
func NewGnomeSecretStore(appId, appDescription string) Secretive {
	//error can be ignored as it has already run successfully in the IsGnome function
	cmd, _ := exec.LookPath(gnomesecretcmd)
	return &GnomeSecretStore{
		cmd:            cmd,
		appIdentifier:  appId,
		appDescription: appDescription,
	}
}

// Store saves a secret value by associating it with a given key in the secret store using the `secret-tool` command.
func (l *GnomeSecretStore) Store(key string, value []byte) error {
	// We have to pipe the output of echo into secret-tool. Echo puts the value on a single line to stdOut.
	// On the systems I have tested, no history is kept.
	cmd1 := exec.Command("echo",
		"-n",
		string(value),
	)
	// This is the actual command we want to run using the piped in key
	cmd2 := exec.Command(
		l.cmd,
		"store",
		fmt.Sprintf(`--label=%s`, l.appDescription),
		l.appIdentifier,
		key,
	)
	cmd2.Stdin, _ = cmd1.StdoutPipe()
	cmd2Output, _ := cmd2.StdoutPipe()
	_ = cmd2.Start()
	_ = cmd1.Start()
	cmd2Result, err := io.ReadAll(cmd2Output)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading command output: %s", string(cmd2Result)))
	}
	_ = cmd1.Wait()
	_ = cmd2.Wait()
	return err
}

// Load retrieves the secret associated with the specified key from the secret store using the `secret-tool` command.
func (l *GnomeSecretStore) Load(key string) ([]byte, error) {
	cmd := exec.Command(
		l.cmd,
		"lookup",
		l.appIdentifier,
		key,
	)
	out, err := cmd.Output()

	// the secret-tool utility will return 1 if the secret doesn't exist
	if cmd.ProcessState.ExitCode() != 0 {
		return []byte{}, ErrSecretNotFound
	}
	return out, err
}

// Exists checks whether a secret associated with the specified key exists in the secret store and returns true if found.
func (l *GnomeSecretStore) Exists(key string) (bool, error) {
	out, err := l.Load(key)
	if err != nil {
		if errors.Is(err, ErrSecretNotFound) {
			return false, nil
		}
		return false, err
	}
	if len(out) > 0 {
		return true, nil
	}
	return false, nil
}

// Delete removes the secret associated with the specified key from the secret store using the `secret-tool` command.
func (l *GnomeSecretStore) Delete(key string) error {
	cmd := exec.Command(
		l.cmd,
		"clear",
		l.appIdentifier,
		key,
	)
	_, err := cmd.Output()
	return err
}

// IsGnome checks if both GNOME Keyring and Secret Tool are installed and returns true if they are, otherwise returns an error.
func IsGnome() (bool, error) {
	_, err := PackageInstalled(gnomekeyringcmd)
	if err != nil {
		return false, ErrKeyringNotInstalled
	}
	_, err = PackageInstalled(gnomesecretcmd)
	if err != nil {
		return false, ErrSecretManagerNotInstalled
	}
	return true, nil
}
