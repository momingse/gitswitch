package cmd

import (
	"gs/libs"

	"github.com/spf13/cobra"
)

func NewRootCommand(dbService *libs.DBService, fileService *libs.FileService) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gs",
		Short: "gitswitch: quick and easy Git project switching",
		Long:  "gitswitch (gs) is a fast and simple CLI tool for switching between your Git projects.",
	}

	rootCmd.AddCommand(NewAddCmd(dbService, fileService))
	return rootCmd
}
