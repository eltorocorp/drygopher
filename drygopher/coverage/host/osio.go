package host

import (
	"io/ioutil"
	"os"

	"github.com/eltorocorp/drygopher/drygopher/coverage/hostiface"
)

// OSIO is a wrapper around common OS and IO functions.
type OSIO struct{}

// MustRemove wrapper
func (OSIO) MustRemove(filename string) {
	err := os.Remove(filename)
	if err == nil || err != nil && os.IsNotExist(err) {
		return
	}
	panic("OSIO.MustRemove encountered an unanticipated error condition: " + err.Error())
}

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

var _ hostiface.OSIOAPI = (*OSIO)(nil)
