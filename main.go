package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wow-log-parser/config"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	setLogLevel()

	log.Info("Genesis - WoW Log Parser and Log Fixer")
	log.Infof("Configuration Information: %+v", config.Config)

	// Array of all the lines in log with timestamps
	loadedData := make([]WoWLog, 0)

	// Note that year is not included in the logs.... I.e. midnight Raid over newyear will fail!
	// Also follow up with how things work on months with double digit...
	var dateFormat = "1/02 15:04:05.000"

	iFile, err := os.Open(config.Config.InputFile)
	if err != nil {
		log.WithFields(log.Fields{
			"message":   err,
			"inputfile": config.Config.InputFile,
		}).Fatal("File Load")
	}
	defer iFile.Close()

	scanner := bufio.NewScanner(iFile)
	for scanner.Scan() {
		inputText := scanner.Text()

		t, err := time.Parse(dateFormat, inputText[0:17])
		if err != nil {
			log.WithFields(log.Fields{
				"message":    err,
				"dateFormat": dateFormat,
				"date":       inputText[0:17],
				"line":       inputText,
			}).Fatal("Parse Time")
		}

		loadedData = append(loadedData, WoWLog{t, inputText[19:], inputText})
	}

	log.WithFields(log.Fields{
		"message":   "file load completed",
		"linecount": len(loadedData),
	}).Info("File Loaded")

	// Output error count
	var previousDate = time.Date(0000, time.January, 01, 01, 0, 0, 0, time.UTC)
	var previousString = ""
	var errorCount = 0
	for c, v := range loadedData {
		if v.Timestamp.Before(previousDate) {
			log.WithFields(log.Fields{
				"beforetime":   previousDate,
				"aftertime":    v.Timestamp,
				"beforestring": previousString,
				"afterstring":  v.Complete,
				"linenumber":   c,
			}).Debug("Date error")

			errorCount++
		} else {
			previousDate = v.Timestamp
			previousString = v.Complete
		}
	}

	// Output error summary
	// exit if ther were no errors
	if errorCount == 0 {
		log.WithFields(log.Fields{
			"message": "there were no errors in file",
		}).Info("Date error")
		os.Exit(0)
	} else {
		log.WithFields(log.Fields{
			"message":    "there were errors in file",
			"errorcount": errorCount,
		}).Info("Date error")
	}

	// Exit if job is done
	if config.Config.CheckOnly == true {
		log.WithFields(log.Fields{
			"message": "check only flag was true, exiting",
		}).Info("Program")
		os.Exit(0)
	}

	oFile, err := os.Create(config.Config.OutputFile)
	if err != nil {
		log.WithFields(log.Fields{
			"message":  "could not create file",
			"filename": config.Config.OutputFile,
		}).Fatal("File output")
	}
	defer oFile.Close()

	// Sort in place, original order kept when equal
	sort.SliceStable(loadedData, func(i, j int) bool {
		return loadedData[i].Timestamp.Before(loadedData[j].Timestamp)
	})

	// Output to file
	for _, value := range loadedData {
		formattedString := fmt.Sprintf("%s  %s\n", value.Timestamp.Format(dateFormat), value.Payload)
		_, err = oFile.WriteString(formattedString)

		if err != nil {
			log.WithFields(log.Fields{
				"message":  "could not write to file",
				"filename": config.Config.OutputFile,
				"err":      err,
			}).Fatal("File output")
		}
	}
}

// WoWLog is the timestamp plus data
// this is all we need for sorting the logs
type WoWLog struct {
	Timestamp time.Time
	Payload   string
	Complete  string
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
