package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type Command struct {
	Program string   `yaml:"program"`
	Args    []string `yaml:"args"`
}

type GithubDependency struct {
	Name             string    `yaml:"name"`
	Repo             string    `yaml:"repo"`
	Commands         []Command `yaml:"commands"`
	InRepoPaths      []string  `yaml:"in_repo_paths"`
	SourceDirectives []string  `yaml:"source_directives"`
}

type Dependencies struct {
	RepoDirectory      string             `yaml:"repo_directory"`
	GithubDependencies []GithubDependency `yaml:"github"`
}

// TODO: write to .dsetuprc
// TODO: push to repo, include instructions

func main() {
	yamlBytes, err := os.ReadFile("test.yaml")

	if err != nil {
		errorExit(err)
	}

	var dependencies Dependencies
	yaml.Unmarshal(yamlBytes, &dependencies)

	home, err := os.UserHomeDir()
	if err != nil {
		errorExit(err)
	}
	repoDirectory := home + "/" + strings.Trim(dependencies.RepoDirectory, "/") + "/"

	for _, dependency := range dependencies.GithubDependencies {
		ensureGithubInstallation(repoDirectory, dependency)
	}
}

func ensureGithubInstallation(repoDirectory string, dependency GithubDependency) {
	fmt.Println("Ensuring installation for", dependency.Name)
	repoPath := repoDirectory + path.Base(dependency.Repo)

	if dirExists(repoPath) {
		fmt.Println("Repository is already cloned.")
	} else {
		fmt.Println("Cloning repository")
		cloneRepo(repoDirectory, dependency)
	}

	for _, path := range dependency.InRepoPaths {
		fmt.Println("export PATH=\"$PATH:" + repoPath + "/" + path + "\"")
	}

	for _, path := range dependency.SourceDirectives {
		fmt.Println("source " + repoPath + "/" + path)
	}

}

func cloneRepo(repoDirectory string, dependency GithubDependency) {
	clone := exec.Command("git", "clone", "git@github.com:"+dependency.Repo+".git")
	clone.Dir = repoDirectory
	clone.Stdout = os.Stdout
	clone.Stderr = os.Stderr
	err := clone.Run()

	if err != nil {
		errorExit(err)
	}
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func errorExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
