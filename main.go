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
		fmt.Println("No config file found; using defaults or env vars.")
	}

	cmd.Execute()
}
