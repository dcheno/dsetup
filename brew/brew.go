package brew

import (
	"log"
	"os"
	"os/exec"

	"github.com/dcheno/dsetup/core"
)

type Dependency struct {
	Formula             string `yaml:"formula"`
	core.DependencyInfo `yaml:",inline"`
}

func (dependency Dependency) RelativeBase(config core.Config) string {
	return ""
}

func (dependency Dependency) EnsureInstallation(config core.Config) {
	command := exec.Command("brew", "install", dependency.Formula)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()

	if err != nil {
		log.Fatal(err)
	}
}
