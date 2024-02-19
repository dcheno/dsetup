package github

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/dcheno/dsetup/core"
	"github.com/dcheno/dsetup/utils"
)

type Dependency struct {
	Repo                string         `yaml:"repo"`
	InstallCommands     []core.Command `yaml:"install_commands"`
	core.DependencyInfo `yaml:",inline"`
}

func (dependency Dependency) RelativeBase(config core.Config) string {
	return config.ReposDirectory + "/" + path.Base(dependency.Repo)
}

func (dependency Dependency) EnsureInstallation(config core.Config) {
	repoPath := dependency.RelativeBase(config)

	if utils.DirExists(repoPath) {
		fmt.Println("Repository is already cloned.")
	} else {
		fmt.Println("Cloning repository")
		cloneRepo(config.ReposDirectory, dependency)
	}

	for _, raw_command := range dependency.InstallCommands {
		command := exec.Command(raw_command.Program, raw_command.Args...)
		command.Dir = repoPath
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Run()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func cloneRepo(repoDirectory string, dependency Dependency) {
	clone := exec.Command("git", "clone", "git@github.com:"+dependency.Repo+".git")
	clone.Dir = repoDirectory
	clone.Stdout = os.Stdout
	clone.Stderr = os.Stderr
	err := clone.Run()

	if err != nil {
		log.Fatal(err)
	}
}
