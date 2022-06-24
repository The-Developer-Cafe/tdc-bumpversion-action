package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	_, err := RunCommand("which", "git")

	if err != nil {
		log.Fatal("'which git' failed, please make sure git is installed!", err)
	}

	_, err = RunCommand("git", "status")

	if err != nil {
		log.Fatal("'git status' failed, please make sure command is being run in a git directory", err)
	}

	latestTag, err := RunCommand("bash", "-c", "git tag -l --sort=-creatordate | head -n 1")
	if err != nil {
		log.Fatal("failed to get the latest git tag", err)
	}

	log.Printf("latest tag: %s\n", latestTag)

	minorVersion, err := ParseTagMajorVersion(strings.TrimSpace(latestTag))
	if err != nil {
		log.Fatal("failed to parse tag", err)
	}

	minorVersion++
	newVersion := fmt.Sprintf("v0.%d.0", minorVersion)

	_, err = RunCommand("bash", "-c", fmt.Sprintf("git tag -a %s -m \"%s\"", newVersion, newVersion))
	if err != nil {
		log.Fatal("failed to create new tag ", minorVersion, err)
	}

	log.Printf("pushing %s tag\n", newVersion)

	_, err = RunCommand("bash", "-c", fmt.Sprintf("git push origin tag %s", newVersion))
	if err != nil {
		log.Fatal("failed to push tag to remote repo", minorVersion, err)
	}

	log.Printf("tagged successfully with version %s\n", newVersion)
}

func ParseTagMajorVersion(tag string) (int, error) {
	if len(tag) == 0 || tag[0] != 'v' {
		return 0, errors.New("invalid tag: " + tag)
	}

	splits := strings.Split(tag[1:], ".")
	if len(splits) != 3 {
		return 0, errors.New("invalid tag version: " + tag)
	}

	num, err := strconv.Atoi(splits[1])
	if err != nil {
		return 0, err
	}

	return num, nil
}

func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(out), nil
}

