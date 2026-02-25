//go:build !windows

package daemon

import (
	"fmt"
	"os"

	godaemon "github.com/sevlyar/go-daemon"
)

// Daemonize detaches the process and returns the child process handle.
// If the returned process is nil, this is the child process.
// If the returned process is non-nil, this is the parent process.
func Daemonize() (*os.Process, error) {
	// no PID file needed
	// we don't want log file, server handles its own logging
	ctx := &godaemon.Context{
		PidFileName: "",
		PidFilePerm: 0,
		LogFileName: "",
		LogFilePerm: 0,
		WorkDir:     "/",
		Umask:       027,
		Args:        os.Args,
		Env:         append(os.Environ(), fmt.Sprintf("%s=1", DaemonEnvVar)),
	}

	child, err := ctx.Reborn()
	if err != nil {
		return nil, fmt.Errorf("failed to daemonize: %w", err)
	}

	return child, nil
}
