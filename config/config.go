package config

import (
	"flag"
	"os"
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var Config = struct {
	InputFile  string
	OutputFile string
	CheckOnly  bool
	LogLevel   string
	Split      bool
	Sort       bool
}{
	"input.txt",
	"output.txt",
	true,
	"info",
	false,
	false,
}

var validLogLevel = map[string]bool{
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
	"trace": true,
}

// Init will initialize default parameters
func init() {
	mesg := "Input filename: INPUT_FILE"
	if v := os.Getenv("INPUT_FILE"); v != "" {
		flag.StringVar(&Config.InputFile, "inputfile", v, mesg)
	} else {
		flag.StringVar(&Config.InputFile, "inputfile", Config.InputFile, mesg)
	}

	mesg = "Output filename: OUTPUT_FILENAME"
	if v := os.Getenv("OUTPUT_FILENAME"); v != "" {
		flag.StringVar(&Config.OutputFile, "outputfile", v, mesg)
	} else {
		flag.StringVar(&Config.OutputFile, "outputfile", Config.OutputFile, mesg)
	}

	mesg = "Check only and output errors: CHECK_ONLY"
	if v := os.Getenv("CHECK_ONLY"); v != "" {
		b, _ := strconv.ParseBool(v)
		flag.BoolVar(&Config.CheckOnly, "checkonly", b, mesg)
	} else {
		flag.BoolVar(&Config.CheckOnly, "checkonly", Config.CheckOnly, mesg)
	}

	mesg = "Set application LogLevel. [trace|debug|info|warn|error] LOG_LEVEL"
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		flag.StringVar(&Config.LogLevel, "loglevel", v, mesg)
	} else {
		flag.StringVar(&Config.LogLevel, "loglevel", Config.LogLevel, mesg)
	}

	mesg = "Split output on mulitple RAIDS: SPLIT"
	if v := os.Getenv("SPLIT"); v != "" {
		b, _ := strconv.ParseBool(v)
		flag.BoolVar(&Config.Split, "split", b, mesg)
	} else {
		flag.BoolVar(&Config.Split, "split", Config.Split, mesg)
	}

	mesg = "Sort data on timestamp: SORT"
	if v := os.Getenv("SORT"); v != "" {
		b, _ := strconv.ParseBool(v)
		flag.BoolVar(&Config.Sort, "sort", b, mesg)
	} else {
		flag.BoolVar(&Config.Sort, "sort", Config.Sort, mesg)
	}

	flag.Parse()

	if validLogLevel[Config.LogLevel] == false {
		log.WithFields(log.Fields{
			"parameter": "logLevel",
			"value":     Config.LogLevel,
			"valid":     reflect.ValueOf(validLogLevel).MapKeys(),
		}).Fatal("Config Error")
	}
}
