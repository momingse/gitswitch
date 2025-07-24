package main

import (
	"fmt"
	"fs/cmd"

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
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	cmd.Execute()
}
