package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dcheno/dsetup/core"
	"github.com/dcheno/dsetup/github"
	"gopkg.in/yaml.v3"
)

type DSetup struct {
	Config       core.Config    `yaml:"config"`
	Dependencies DependencyList `yaml:"dependencies"`
}

type DependencyList []core.Dependency

// TODO: push to repo, include instructions
// TODO: ensure existence of repos directory
// TODO: install commands

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS] config-filename\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage

	groupList := core.GroupList{"default"}
	flag.Var(&groupList, "group", "Group to include in setup. List the flag multiple times for multiple groups. 'default' group is always included.")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("Unexpected input. Must take one and only one config file name as argument.")
	}

	configFilename := flag.Arg(0)

	yamlBytes, err := os.ReadFile(configFilename)

	if err != nil {
		log.Fatal(err)
	}

	var dsetup DSetup
	yaml.Unmarshal(yamlBytes, &dsetup)

	config := dsetup.Config

	assertString(config.ReposDirectory)
	config.ReposDirectory = strings.TrimRight(os.ExpandEnv(config.ReposDirectory), "/")

	assertString(config.EnvFilename)
	assertString(config.RcFilename)

	if config.FishFilename != "" {
		createOrTruncateAutoGeneratedFile(config.FishFilename)
	}

	if config.EnvFilename != "" {
		createOrTruncateAutoGeneratedFile(config.EnvFilename)
	}

	if config.RcFilename != "" {
		createOrTruncateAutoGeneratedFile(config.RcFilename)
	}

	for _, dependency := range dsetup.Dependencies {
		if dependency.HasAtLeastOneGroup(groupList) {
			if dependency.CommandExists() {
				fmt.Println(dependency.Name(), "already installed.")
				continue
			}
			fmt.Println("Ensuring installation for", dependency.Name())
			dependency.EnsureInstallation(config)
			dependency.WriteFiles(config, dependency.RelativeBase(config))
		}
	}
}

func assertString(s string) {
	if s == "" {
		log.Fatal(s + " must be provided.")
	}
}

func createOrTruncateAutoGeneratedFile(filename string) {
	f, err := os.Create(os.ExpandEnv(filename))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(f, "# *******     AUTOGENERATED FILE     *******")
	fmt.Fprintln(f, "# ----- created by dsetup. DO NOT EDIT -----")
	fmt.Fprintln(f, "")
}

type typed struct {
	Type string `yaml:"type"`
}

func (dependencies *DependencyList) UnmarshalYAML(value *yaml.Node) error {
	for _, rawDependency := range value.Content {
		var t typed
		err := rawDependency.Decode(&t)
		if err != nil {
			return err
		}

		var decoded core.Dependency
		switch t.Type {
		case "github":
			decoded = new(github.Dependency)
		}

		err = rawDependency.Decode(decoded)
		if err != nil {
			return err
		}
		*dependencies = append(*dependencies, decoded)

	}
	return nil
}
