package hostiface

import "os"

// OSIOAPI represents something that provides OS and IO methods.
type OSIOAPI interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	ReadDir(dirname string) ([]os.FileInfo, error)
	LookupEnv(key string) (string, bool)
}
