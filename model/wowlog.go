package model

import (
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const logHeader = "COMBAT_LOG_VERSION"

// WoWLogFile Contains raw data and logs from raid
type WoWLogFile struct {
	OriginalData []string
	LogEntries   []WoWLogEntry
}

// WoWLogEntry is represents a single line
type WoWLogEntry struct {
	Timestamp time.Time
	Payload   string
	Complete  string
}

// Split will return an array of logfiles
func (r *WoWLogFile) Split() *[]WoWLogFile {
	wowlogs := make([]WoWLogFile, r.CountRaids())
	raidCount := -1

	for _, v := range r.LogEntries {
		if strings.HasPrefix(v.Payload, logHeader) {
			raidCount++
		}
		wowlogs[raidCount].LogEntries = append(wowlogs[raidCount].LogEntries, v)
	}
	return &wowlogs
}

// PrintInfo will output stuff to logs
func (r *WoWLogFile) PrintInfo() {
	if len(r.LogEntries) == 0 {
		log.WithFields(log.Fields{
			"message": "Log has no antries",
		}).Fatal("Program")
	}

	log.WithFields(log.Fields{
		"firstentrytime": r.LogEntries[0].Timestamp,
		"lastentrytime":  r.LogEntries[len(r.LogEntries)-1].Timestamp,
	}).Info("Status")
}

// SortOrder will sort according to datetime order
func (r *WoWLogFile) SortOrder() {
	// Sort in place, original order kept when equal
	sort.SliceStable(r.LogEntries, func(i, j int) bool {
		return r.LogEntries[i].Timestamp.Before(r.LogEntries[j].Timestamp)
	})
}

// CountRaids returns number of raids in log
func (r *WoWLogFile) CountRaids() int {
	numberOfRaids := 0

	for _, v := range r.LogEntries {
		if strings.HasPrefix(v.Payload, logHeader) {
			numberOfRaids++
		}
	}
	return numberOfRaids
}

// DateErrorCount will check order
func (r *WoWLogFile) DateErrorCount() int {
	// Output error count
	var previousDate = time.Date(0000, time.January, 01, 01, 0, 0, 0, time.UTC)
	var previousString = ""
	var errorCount = 0
	for c, v := range r.LogEntries {
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
	return errorCount
}
