package testdir

import (
	"os"
)

// New creates a working directory at the specified path.
// A cleanup function is returned to delete all content.
func New(p string) (cleanup func() error, err error) {
	err = os.Mkdir(p, os.ModePerm)
	if err != nil {
		return func() error { return nil }, err
	}
	return func() error {
		return os.RemoveAll(p)
	}, nil
}
