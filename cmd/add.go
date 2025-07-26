//go:generate mockgen -destination=../mocks/cmd/add.go -package=mocks -source=add.go
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

type DBService interface {
	Add(alias string, path string) error
}

type FileService interface {
	GetCurrentPath() (string, error)
	GetParentFolderName(path string) string
	CheckIfPathExists(path string) (bool, error)
}

func NewAddCmd(dbService DBService, fileService FileService) *cobra.Command {
	return &cobra.Command{
		Use:   "add [alias] [path]",
		Short: "Add a Git project to gitswitch with an alias",
		Long: `Add a Git project to gitswitch so you can quickly switch to it later using 'gs <alias>'.

Usage Scenarios:
  1. gs add                → Adds the current directory using its folder name as alias (if it's a Git repo)
  2. gs add <alias>        → Adds the current directory with a custom alias
  3. gs add <alias> <path> → Adds the specified path with the given alias

In all cases, the path is saved and can be accessed later with 'gs <alias>'.`,
		Args: cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias, path, err := determineAliasAndPath(args, fileService)
			if err != nil {
				return err
			}

			if err := dbService.Add(alias, path); err != nil {
				return fmt.Errorf("failed to add %s with alias %s", path, alias)
			}

			return nil
		},
	}
}

func determineAliasAndPath(args []string, fileService FileService) (string, string, error) {
	currentPath, err := getPath(args, fileService)
	if err != nil {
		return "", "", err
	}

	alias, err := getAlias(args, fileService, currentPath)
	if err != nil {
		return "", "", err
	}

	return alias, currentPath, nil
}

func getPath(args []string, fileService FileService) (string, error) {
	switch len(args) {
	case 0, 1:
		currentPath, err := fileService.GetCurrentPath()
		if err != nil {
			return "", errors.New("failed to get current path")
		}
		return currentPath, nil
	case 2:
		isPathExists, err := fileService.CheckIfPathExists(args[1])
		if err != nil || !isPathExists {
			return "", errors.New("path does not exist")
		}
		return args[1], nil
	}
	return "", errors.New("invalid number of arguments")
}

func getAlias(args []string, fileService FileService, currentPath string) (string, error) {
	switch len(args) {
	case 0:
		return fileService.GetParentFolderName(currentPath), nil
	case 1, 2:
		return args[0], nil
	}
	return "", errors.New("invalid number of arguments")
}
