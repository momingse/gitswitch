package cmd

import (
	"gs/libs"

	"github.com/spf13/cobra"
)

func NewAddCmd(dbService *libs.DBService, fileService *libs.FileService) *cobra.Command {
	return &cobra.Command{
		Use:   "add [alias] [path]",
		Short: "Add a new project alias to gitswitch",
		Long:  "Registers a new alias and its corresponding path into the gitswitch database.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
		}}
}
