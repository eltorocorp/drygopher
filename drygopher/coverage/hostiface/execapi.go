package hostiface

// ExecAPI represents anything that knows how to build a command.
type ExecAPI interface {
	Command(name string, arg ...string) CommandAPI
}

// CommandAPI represents a command.
type CommandAPI interface {
	Run() error
	CombinedOutput() ([]byte, error)
	Output() ([]byte, error)
}
