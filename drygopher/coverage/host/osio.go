package host

import (
	"io/ioutil"
	"os"
)

// OSIO is a wrapper around common OS and IO functions.
type OSIO struct{}

// ReadFile wrapper
func (OSIO) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// WriteFile wrapper
func (OSIO) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

// ReadDir wrapper
func (OSIO) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

//LookupEnv wrapper
func (OSIO) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
