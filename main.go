package main

import (
	"fmt"
	"gs/cmd"
	"gs/libs"
	"os"

	"github.com/spf13/viper"
)

func errorHandler(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error %s: %v\n", msg, err)
		os.Exit(1)
	}
}

func main() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			errorHandler(err, "Config file not found")
		} else {
			errorHandler(err, "Config file error")
		}
	}

	bucketName := viper.GetString("kv_bucket_name")
	db, err := libs.NewBoltDB(bucketName)
	if err != nil {
		errorHandler(err, "NewBoltDB error")
	}

	dbService := libs.NewDBService(db, bucketName)

	rootCmd := cmd.NewRootCommand(dbService)
	if err := rootCmd.Execute(); err != nil {
		errorHandler(err, "Execute error")
	}
}
