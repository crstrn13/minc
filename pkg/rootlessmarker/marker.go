package rootlessmarker

import (
	"os"
	"path/filepath"
)

const markerPrefix = ".rootless-cluster"

// Path returns the path to the on-disk marker for the named cluster, or an error
// if the user config dir cannot be resolved.
func Path(name string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "minc", markerPrefix+"-"+name), nil
}

// Present is true when the named minc cluster was created with rootless Podman
// (see register).
func Present(name string) bool {
	p, err := Path(name)
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}

// Set writes the marker after a successful rootless create.
func Set(name string) error {
	p, err := Path(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	return os.WriteFile(p, []byte("1\n"), 0644)
}

// Remove deletes the marker; call after a successful minc delete of a rootless cluster.
func Remove(name string) error {
	p, err := Path(name)
	if err != nil {
		return err
	}
	err = os.Remove(p)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
