package custom

import (
	"log"
	"os"
	"os/exec"

	"github.com/dcheno/dsetup/core"
)

type Dependency struct {
	InstallCommands     []core.Command `yaml:"install_commands"`
	core.DependencyInfo `yaml:",inline"`
}

func (dependency Dependency) RelativeBase(config core.Config) string {
	return ""
}

func (dependency Dependency) EnsureInstallation(config core.Config) {
	for _, raw_command := range dependency.InstallCommands {
		command := exec.Command(raw_command.Program, raw_command.Args...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Run()

		if err != nil {
			log.Fatal(err)
		}
	}
}
