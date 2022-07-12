package cli

import "os/exec"

type CommandExecutor interface {
	Execute(string, []string) (string, error)
}

type CLICommandExecutor struct{}

func (CLICommandExecutor) Execute(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.Output()

	return string(stdout), err
}

func NewCommandExecutor() CommandExecutor {
	return CLICommandExecutor{}
}
