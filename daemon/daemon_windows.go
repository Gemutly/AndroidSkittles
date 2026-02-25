//go:build windows

package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Daemonize spawns the current process as a detached background process on Windows.
// Returns the child process handle (non-nil = parent, nil = child).
func Daemonize() (*os.Process, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	cmd := exec.Command(executable, os.Args[1:]...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=1", DaemonEnvVar))

	// detach from the parent console so the child survives parent exit
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	// discard stdout/stderr; the server handles its own logging
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to daemonize: %w", err)
	}

	// detach so the child is not reaped when this process exits
	_ = cmd.Process.Release()

	return cmd.Process, nil
}
