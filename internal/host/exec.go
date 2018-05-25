package host

import (
	"os/exec"

	"github.com/eltorocorp/drygopher/internal/hostiface"
)

// Exec is a wrapper around common exec functions.
type Exec struct{}

// Command wrapper
func (Exec) Command(name string, arg ...string) hostiface.CommandAPI {
	return exec.Command(name, arg...)
}
