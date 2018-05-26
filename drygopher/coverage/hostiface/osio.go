package hostiface

import "os"

// OSIOAPI represents something that provides OS and IO methods.
type OSIOAPI interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	ReadDir(dirname string) ([]os.FileInfo, error)
	LookupEnv(key string) (string, bool)
}

// FileInfo is a wrapper around the os.FileInfo interface.
// Mockery isn't very good at making mocks for interfaces
// outside of the project, so I'm wrapping the FileInfo
// interface to make life a bit easier.
type FileInfo interface {
	os.FileInfo
}
