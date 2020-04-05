package controller

import (
	"bufio"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wow-log-parser/config"
	"github.com/wow-log-parser/model"
)

const dateFormat = "1/02 15:04:05.000"

// LoadFile will load a file
func LoadFile(fileName string) *model.WoWLogFile {
	wowlogfile := new(model.WoWLogFile)

	iFile, err := os.Open(fileName)
	if err != nil {
		log.WithFields(log.Fields{
			"message":   err,
			"inputfile": fileName,
		}).Fatal("File Load")
	}
	defer iFile.Close()

	scanner := bufio.NewScanner(iFile)
	for scanner.Scan() {
		inputText := scanner.Text()

		// Check the line contains a propper date,
		// as a simple check
		t, err := time.Parse(dateFormat, inputText[0:17])
		if err != nil {
			log.WithFields(log.Fields{
				"message":    err,
				"dateFormat": dateFormat,
				"date":       inputText[0:17],
				"line":       inputText,
			}).Fatal("Parse Time")
		}

		wowlogfile.OriginalData = append(wowlogfile.OriginalData, inputText)
		wowlogfile.LogEntries = append(wowlogfile.LogEntries,
			model.WoWLogEntry{
				Timestamp: t,
				Payload:   inputText[19:],
				Complete:  inputText,
			})
	}

	log.WithFields(log.Fields{
		"message":   "file loaded",
		"inputfile": fileName,
		"lines":     len(wowlogfile.OriginalData),
	}).Info("Status")

	return wowlogfile
}

// SaveFile will save stuff to file
func SaveFile(fileName string, wowlog *model.WoWLogFile) {
	oFile, err := os.Create(fileName)
	if err != nil {
		log.WithFields(log.Fields{
			"message":  "could not create file",
			"filename": fileName,
		}).Fatal("File output")
	}
	defer oFile.Close()

	// Output to file
	for _, value := range wowlog.LogEntries {
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
