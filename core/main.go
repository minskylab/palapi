package main

import (
	"github.com/minskylab/palapi"
	rae "github.com/minskylab/palapi/providers/RAE"
	wordreference "github.com/minskylab/palapi/providers/WordReference"
	"github.com/spf13/viper"
)

const DefaultDatabaseFilename = "palapi.db"
const DatabaseFilenameKey = "DatabaseFilename"

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".palapi")

	viper.SetDefault(DatabaseFilenameKey, DefaultDatabaseFilename)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(err)
		}
	}

	storage, err := palapi.NewDefaultPersistence(viper.GetString(DatabaseFilenameKey))
	if err != nil {
		panic(err)
	}

	manager, err := palapi.NewManager(storage, 2)
	if err != nil {
		panic(err)
	}

	raeProvider, err := rae.NewProvider()
	if err != nil {
		panic(err)
	}

	wordReferenceProvider, err := wordreference.NewProvider()
	if err != nil {
		panic(err)
	}

	if err := manager.RegisterProvider(raeProvider); err != nil {
		panic(err)
	}

	if err := manager.RegisterProvider(wordReferenceProvider); err != nil {
		panic(err)
	}

}
