package main

import (
	"fmt"
	"fs/cmd"
	"os"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found; using defaults or env vars.")
		} else {
			fmt.Fprintf(os.Stderr, "fatal error config file: %v\n", err)
			os.Exit(1)
		}
	}

	cmd.Execute()
}
