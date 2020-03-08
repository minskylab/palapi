package main

import (
	"github.com/k0kubun/pp"
	"github.com/minskylab/palapi"
	rae "github.com/minskylab/palapi/providers/rae"
	wordreference "github.com/minskylab/palapi/providers/wordreference"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const DefaultDatabaseFilename = "palapi.db"
const DatabaseFilenameKey = "DatabaseFilename"

const DebugKey = "Debug"

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(".palapi")

	viper.SetDefault(DatabaseFilenameKey, DefaultDatabaseFilename)
	viper.SetDefault(DebugKey, false)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(err)
		}
	}

	if viper.GetBool(DebugKey) {
		log.SetLevel(log.DebugLevel)
	}

	storage, err := palapi.NewDefaultPersistence(viper.GetString(DatabaseFilenameKey))
	if err != nil {
		panic(err)
	}

	manager, err := palapi.NewManager(storage, 1)
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

	gov, err := manager.ReportWord("desdeñoso")
	if err != nil {
		panic(errors.Cause(err))
	}

	pp.Println(gov)

	// word, err := storage.GetWord("gobernar")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// pp.Println(word)

}
