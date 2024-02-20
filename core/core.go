package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dcheno/dsetup/utils"
)

type Config struct {
	ReposDirectory string `yaml:"repos_directory"`
	EnvFilename    string `yaml:"env_file"`
	RcFilename     string `yaml:"rc_file"`
	FishFilename   string `yaml:"fish_file"`
}

type Dependency interface {
	Name() string
	EnsureInstallation(config Config)
	Exists() bool
	HasAtLeastOneGroup(checkGroups GroupList) bool
}

type FileWriter interface {
	WriteFiles(config Config, relativeBase string)
	RelativeBase(config Config) string
}

type Command struct {
	Program string   `yaml:"program"`
	Args    []string `yaml:"args"`
}

type DotFile struct {
	AbsoluteSourceDirectives []string `yaml:"absolute_source_directives"`
	RelativeSourceDirectives []string `yaml:"relative_source_directives"`
	AbsolutePaths            []string `yaml:"absolute_paths"`
	RelativePaths            []string `yaml:"relative_paths"`
}

type FishFile struct {
	AbsolutePaths []string `yaml:"absolute_paths"`
	RelativePaths []string `yaml:"relative_paths"`
}

type DependencyInfo struct {
	Env     DotFile   `yaml:".env"`
	Rc      DotFile   `yaml:".rc"`
	Fish    FishFile  `yaml:".fish"`
	Command string    `yaml:"command"`
	Groups  GroupList `yaml:"groups"`
}

type GroupList []string

func (gl *GroupList) String() string {
	if gl != nil {
		strings.Join(*gl, ",")
	}
	return ""
}

func (gl *GroupList) Set(group string) error {
	*gl = append(*gl, group)
	return nil
}

func (dependencyInfo DependencyInfo) Exists() bool {
	return utils.CommandExists(dependencyInfo.Command)
}

func (dependencyInfo DependencyInfo) Name() string {
	return dependencyInfo.Command
}

func (dependencyInfo DependencyInfo) HasAtLeastOneGroup(checkGroups GroupList) bool {
	return HasAtLeastOneGroup(dependencyInfo.Name(), dependencyInfo.Groups, checkGroups)
}

func HasAtLeastOneGroup(name string, dependencyGroups GroupList, checkGroups GroupList) bool {
	if len(dependencyGroups) == 0 {
		fmt.Printf("⚠️ '%s' is not attached to any groups and will not be installed! ⚠️⚠\n", name)
	}

	for _, includedGroup := range checkGroups {
		for _, dependencyGroup := range dependencyGroups {
			if dependencyGroup == includedGroup {
				return true
			}
		}
	}
	return false
}

func (dependencyInfo DependencyInfo) WriteFiles(config Config, relativeBase string) {
	if config.FishFilename != "" {
		appendFishFile(config.FishFilename, dependencyInfo.Fish, relativeBase)
	}

	if config.EnvFilename != "" {
		appendDotFile(config.EnvFilename, dependencyInfo.Env, relativeBase)
	}

	if config.RcFilename != "" {
		appendDotFile(config.RcFilename, dependencyInfo.Rc, relativeBase)
	}
}

func appendDotFile(filename string, dotfile DotFile, relativeBase string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range dotfile.RelativePaths {
		fmt.Fprintln(f, "export PATH=\"$PATH:"+relativeBase+"/"+path+"\"")
	}

	for _, path := range dotfile.AbsolutePaths {
		fmt.Fprintln(f, "export PATH=\"$PATH:"+path+"\"")
	}

	for _, path := range dotfile.RelativeSourceDirectives {
		fmt.Fprintln(f, "source "+relativeBase+"/"+path)
	}

	for _, path := range dotfile.AbsoluteSourceDirectives {
		fmt.Fprintln(f, "source "+path)
	}
}

func appendFishFile(filename string, fishFile FishFile, relativeBase string) {
	f, err := os.OpenFile(os.ExpandEnv(filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range fishFile.RelativePaths {
		fmt.Fprintln(f, "fish_add_path "+relativeBase+"/"+path)
	}

	for _, path := range fishFile.AbsolutePaths {
		fmt.Fprintln(f, "fish_add_path "+path)
	}
}
