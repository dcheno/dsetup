package directory

import (
	"io/fs"
	"log"
	"os"

	"github.com/dcheno/dsetup/core"
)

type Dependency struct {
	Path        string         `yaml:"path"`
	Permissions fs.FileMode    `yaml:"permissions"`
	Groups      core.GroupList `yaml:"groups"`
}

func (dependency Dependency) Name() string {
	return dependency.Path
}

func (dependency Dependency) EnsureInstallation(_ core.Config) {
	err := os.MkdirAll(os.ExpandEnv(dependency.Path), dependency.Permissions)
	if err != nil {
		log.Fatal(err)
	}
}

func (dependency Dependency) Exists() bool {
	_, err := os.Stat(os.ExpandEnv(dependency.Path))

	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Fatal(err)
		// should be unreachable
		return false
	}
}

func (dependency Dependency) HasAtLeastOneGroup(checkGroups core.GroupList) bool {
	return core.HasAtLeastOneGroup(dependency.Name(), dependency.Groups, checkGroups)
}
