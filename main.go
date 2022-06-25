package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v2"
)

type IncrementConfig struct {
	GitName       string `validate:"min=1"`
	GitEmail      string `validate:"email"`
	IncrementType string `validate:"oneof=major minor patch"`
}

func main() {
	var (
		gitName       string
		gitEmail      string
		incrementType string
	)

	app := &cli.App{
		Name:  "bumpversion",
		Usage: "Increment semantic version and make a git tag",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "gitEmail",
				Required:    true,
				Usage:       "Email to configure git",
				Destination: &gitEmail,
			},
			&cli.StringFlag{
				Name:        "gitName",
				Required:    true,
				Usage:       "Name to configure git",
				Destination: &gitName,
			},
			&cli.StringFlag{
				Name:        "incrementType",
				Required:    true,
				Usage:       "Type of version to increment",
				Destination: &incrementType,
			},
		},
		Action: func(ctx *cli.Context) error {
			cfg := IncrementConfig{
				GitName:       gitName,
				GitEmail:      gitEmail,
				IncrementType: incrementType,
			}

			log.Println("\n", cfg.String())

			err := cfg.Validate()
			if err != nil {
				return err
			}

			return IncrementGitTag(cfg)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("failed to bumpversion: %v\n", err)
	}
}

func IncrementGitTag(cfg IncrementConfig) error {
	pwd, err := RunCommand("pwd")
	if err != nil {
		return errors.New(fmt.Sprint("failed to get current working directory ", err))
	}
	log.Printf("running in folder %s", pwd)

	ls, err := RunCommand("ls", "-la")
	if err != nil {
		return errors.New(fmt.Sprint("failed to get current directory contents", err))
	}
	log.Printf("contents of directory: \n%s", ls)

	_, err = RunCommand("which", "git")

	if err != nil {
		return errors.New(fmt.Sprint("'which git' failed, please make sure git is installed!", err))
	}

	_, err = RunCommand("git", "status")
	if err != nil {
		return errors.New(fmt.Sprint("'git status' failed, please make sure command is being run in a git directory", err))
	}

	_, err = RunCommand("git", "fetch", "--tags")
	if err != nil {
		return fmt.Errorf("failed to fetch git tags: %v", err)
	}

	// Configure git user name and email.
	_, err = RunCommand("git", "config", "--global", "user.name", fmt.Sprintf("'%s'", cfg.GitName))
	if err != nil {
		return errors.New(fmt.Sprint("failed to configure git user.name: ", err))
	}

	_, err = RunCommand("git", "config", "--global", "user.email", fmt.Sprintf("'%s'", cfg.GitEmail))
	if err != nil {
		return errors.New(fmt.Sprint("failed to configure git user.email: ", err))
	}

	latestTag, err := RunCommand("sh", "-c", "git tag -l --sort=-creatordate | head -n 1")
	if err != nil {
		return errors.New(fmt.Sprint("failed to get the latest git tag: ", err))
	}

	log.Printf("latest tag: %s", latestTag)

	newVersion, err := IncrementVersion(latestTag, cfg.IncrementType)
	if err != nil {
		return fmt.Errorf("failed to increment tag version: %e", err)
	}

	log.Printf("bumped tag: %s\n", newVersion)

	_, err = RunCommand("sh", "-c", fmt.Sprintf("git tag -a %s -m \"%s\"", newVersion, newVersion))
	if err != nil {
		return errors.New(fmt.Sprint("failed to create new tag ", newVersion, err))
	}

	log.Printf("pushing %s tag\n", newVersion)

	_, err = RunCommand("sh", "-c", fmt.Sprintf("git push origin tag %s", newVersion))
	if err != nil {
		return fmt.Errorf("failed to push tag to remote repo: %v", err)
	}

	return fmt.Errorf("tagged successfully with version %s", newVersion)
}

func IncrementVersion(existingVersion, incrementType string) (string, error) {
	parsedVersion, err := ParseSemanticVersion(existingVersion)
	if err != nil {
		return "", fmt.Errorf("error parsing existing tag version: %v", err)
	}

	switch incrementType {
	case "major":
		parsedVersion[0]++
		parsedVersion[1] = 0
		parsedVersion[2] = 0
	case "minor":
		parsedVersion[1]++
		parsedVersion[2] = 0
	case "patch":
		parsedVersion[2]++
	default:
		return "", fmt.Errorf("invalid incrementType %s", incrementType)
	}

	return fmt.Sprintf("v%d.%d.%d", parsedVersion[0], parsedVersion[1], parsedVersion[2]), nil
}

func ParseSemanticVersion(tag string) ([]int, error) {
	if len(tag) == 0 || tag[0] != 'v' {
		return nil, errors.New("invalid tag: " + tag)
	}

	splits := strings.Split(strings.Trim(tag[1:], "\n"), ".")
	if len(splits) != 3 {
		return nil, errors.New("invalid tag version: " + tag)
	}

	major, err := strconv.Atoi(splits[0])
	if err != nil {
		return nil, err
	}

	minor, err := strconv.Atoi(splits[1])
	if err != nil {
		return nil, err
	}

	patch, err := strconv.Atoi(splits[2])
	if err != nil {
		return nil, err
	}

	return []int{major, minor, patch}, nil
}

func RunCommand(name string, args ...string) (string, error) {
	log.Println("Running:", name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)

	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (cfg *IncrementConfig) Validate() error {
	err := validator.New().Struct(cfg)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			errStr := strings.Builder{}

			for _, err := range validationErrors {
				errStr.WriteString(err.Error())
				errStr.WriteString("\n")
			}

			return errors.New(errStr.String())
		}

		return err
	}

	return nil
}

func (cfg *IncrementConfig) String() string {
	return fmt.Sprintf("IncrementConfig:\n\tGitEmail: %s\n\tGitUser: %s\n\tIncrementType: %s\n", cfg.GitEmail, cfg.GitName, cfg.IncrementType)
}
