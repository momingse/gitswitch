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
	CheckIfPathExists(path string) bool
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
			switch len(args) {
			case 0:
				currentDirectory, err := fileService.GetCurrentPath()
				if err != nil {
					fmt.Println(err)
					return errors.New("failed to get current directory")
				}

				folderName := fileService.GetParentFolderName(currentDirectory)

				err = dbService.Add(folderName, currentDirectory)
				if err != nil {
					fmt.Println(err)
					return errors.New("failed to add to database")
				}

				fmt.Printf("Added %s to gitswitch\n", folderName)
			case 1:
				alias := args[0]
				currentDirectory, err := fileService.GetCurrentPath()
				if err != nil {
					fmt.Println(err)
					return errors.New("failed to get current directory")
				}

				err = dbService.Add(alias, currentDirectory)
				if err != nil {
					fmt.Println(err)
					return errors.New("failed to add to database")
				}

				fmt.Printf("Added %s to gitswitch\n", alias)
			case 2:
				alias := args[0]
				path := args[1]

				isPathExist := fileService.CheckIfPathExists(path)
				if !isPathExist {
					return errors.New("path does not exist")
				}

				err := dbService.Add(alias, path)
				if err != nil {
					fmt.Println(err)
					return errors.New("failed to add to database")
				}

				fmt.Printf("Adding %s with path %s to gitswitch\n", alias, path)
			default:
				return errors.New("unexpected number of arguments")
			}

			return nil
		},
	}
}
