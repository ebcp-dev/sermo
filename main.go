package main

import (
	"log"
	"os"

	"github.com/ebcp-dev/sermo/app"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	// Log current environment.
	if os.Getenv("ENV") == "" {
		os.Setenv("ENV", "dev")
	}
	log.Println("ENV: " + os.Getenv("ENV"))

	a := app.App{}

	a.Initialize()
	if os.Getenv("PORT") == "" {
		// Get port from config if no env variable.
		a.Run(":" + viper.GetString("PORT"))
	} else {
		// Get port from env.
		a.Run(":" + os.Getenv("PORT"))
	}
}
