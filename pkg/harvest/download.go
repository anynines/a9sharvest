package harvest

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/gogolok/go-harvest/harvest"
	log "github.com/sirupsen/logrus"
)

func DownloadAndOutputEntries(verboseFlag bool, outputFilename string) error {
	if verboseFlag {
		log.SetLevel(log.DebugLevel)
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) > 0 {
		lvl, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Fatalf("Wrong log level %v", logLevel)
		}
		log.SetLevel(lvl)
	}

	err := CheckRequiredEnvVariables()
	if err != nil {
		return err
	}

	entries, err := fetchTimeEntries()
	if err != nil {
		return err
	}

	var writer *bufio.Writer
	if len(outputFilename) < 1 {
		writer = bufio.NewWriter(os.Stdout)
	} else {
		file, err := os.Create(outputFilename)
		if err != nil {
			log.Errorf("Failed to open file %v\n", outputFilename)
			return err
		}
		writer = bufio.NewWriter(file)
	}
	defer writer.Flush()

	err = outputTimeEntries(entries, writer)
	if err != nil {
		return err
	}

	return nil
}

func outputTimeEntries(entries []*harvest.TimeEntry, writer io.Writer) error {
	jsonData, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	_, err = writer.Write(jsonData)
	return err
}
