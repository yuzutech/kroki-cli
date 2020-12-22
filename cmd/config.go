package cmd

import "github.com/spf13/viper"

func SetupConfig() {
	// Default values
	viper.SetDefault("endpoint", "https://demo.kroki.io")
	viper.SetDefault("timeout", "20s")

	// Config file name
	viper.SetConfigName("kroki")

	// Default config paths
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("kroki")
	err := viper.BindEnv("endpoint")
	if err != nil {
		exit(err)
	}
	err = viper.BindEnv("timeout")
	if err != nil {
		exit(err)
	}
}


func InitDefaultConfig() {
	_ = viper.ReadInConfig() // ignore error
}