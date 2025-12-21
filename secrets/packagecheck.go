package secrets

import "os/exec"

// PackageInstalled checks if the specified executable is available in the system's PATH and returns its availability status.
func PackageInstalled(pkg string) (bool, error) {
	_, err := exec.LookPath(pkg)

	// check error
	if err != nil {
		// the executable is not found, return false
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return false, nil
		}
		// another kind of error happened
		return false, err
	}

	return true, nil
}
