package utils

import (
	"os"
	"os/exec"
)

func DirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
