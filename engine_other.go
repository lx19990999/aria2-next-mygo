//go:build !windows

package main

import "os/exec"

func configureEngineCmd(_ *exec.Cmd) {
	// Unix: leave process group as-is so a detached child can outlive the GUI
	// when "stop on exit" is disabled.
}
