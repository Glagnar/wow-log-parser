package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/wow-log-parser/config"
	"github.com/wow-log-parser/controller"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	setLogLevel()

	log.Info("Genesis - WoW Log Parser, splitter and Log Fixer")
	log.Infof("Configuration Information: %+v", config.Config)

	wowlog := controller.LoadFile(config.Config.InputFile)

	log.WithFields(log.Fields{
		"message":   "logfile parsed",
		"raidcount": wowlog.CountRaids(),
		"linecount": len(wowlog.LogEntries),
	}).Info("Status")

	log.WithFields(log.Fields{
		"message":    "number of our of order date errors in file",
		"errorcount": wowlog.DateErrorCount(),
	}).Info("Status")

	// Exit if job is done
	if config.Config.CheckOnly == true {
		log.WithFields(log.Fields{
			"message": "check only flag was true, exiting",
		}).Info("Program")
		os.Exit(0)
	}

	if config.Config.Sort == true {
		wowlog.SortOrder()
	}

	if config.Config.Split == true {
		splitRaids := wowlog.Split()

		for c, v := range *splitRaids {
			v.PrintInfo()
			controller.SaveFile(fmt.Sprintf("%s.%d", config.Config.OutputFile, c), &v)
		}
	} else {
		controller.SaveFile(config.Config.OutputFile, wowlog)
	}

}

func setLogLevel() {
	switch config.Config.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
