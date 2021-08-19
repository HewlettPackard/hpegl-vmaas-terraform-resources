// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"os"

	"github.com/spf13/viper"
)

func ReadAccConfig(path string) {
	if os.Getenv("TF_ACC") == "true" {
		viper.AddConfigPath(path)
		viper.SetConfigType("yaml")
		viper.SetConfigName(os.Getenv("TF_ACC_CONFIG"))
		err := viper.ReadInConfig()
		if err != nil {
			panic("fatal error config file: " + err.Error())
		}
	}
}
