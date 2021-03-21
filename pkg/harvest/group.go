package harvest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gogolok/go-harvest/harvest"
	log "github.com/sirupsen/logrus"
)

type OutputFormat uint32

type Result struct {
	GroupedByTag map[string]float64
}

func Group(verboseFlag bool, inputFlag string, outputFlag string, skipUnknown bool) error {
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

	err = CheckOptionalEnvVariables()
	if err != nil {
		return err
	}

	outputFormat, err := parseOutputFormat(outputFlag)
	if err != nil {
		log.Fatalln(err)
	}

	var entries []*harvest.TimeEntry
	if len(inputFlag) < 1 {
		entries, err = fetchTimeEntries()
		if err != nil {
			return err
		}
	} else {
		data, err := os.ReadFile(inputFlag)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, &entries); err != nil {
			return err
		}
	}

	var matcher EntryMatcher
	if len(os.Getenv("PATTERN")) > 0 {
		matcher, err = NewPatternEntryMatcher(os.Getenv("PATTERN"))
		if err != nil {
			return err
		}
	} else {
		matcher = NewTagEntryMatcher(os.Getenv("TAGS"))
	}

	report := NewReport(entries, matcher, skipUnknown)
	report.Run()
	stats := report.Stats()

	var formatter ReportInterface
	switch outputFormat {
	case CSVOutputFormat:
		formatter = NewCSVFormatter(stats)
	default:
		formatter = NewTextFormatter(stats)
	}
	formatter.Output()

	return nil
}

func CheckRequiredEnvVariables() error {
	keys := []string{"ACCOUNT_ID", "TOKEN"}

	for _, key := range keys {
		if len(os.Getenv(key)) < 1 {
			return fmt.Errorf("You MUST set the environment variable %s", key)
		}
	}

	return nil
}

func CheckOptionalEnvVariables() error {
	tagsExist := len(os.Getenv("TAGS")) > 0
	patternExists := len(os.Getenv("PATTERN")) > 0
	if !tagsExist && !patternExists {
		return fmt.Errorf("You MUST set the environment variable TAGS or PATTERN")
	}

	return nil
}

func fetchTimeEntries() ([]*harvest.TimeEntry, error) {
	to := os.Getenv("TO")
	if len(to) < 1 {
		to = time.Now().Format("20060102")
	}
	from := os.Getenv("FROM")
	if len(from) < 1 {
		from = time.Now().AddDate(0, 0, -14).Format("20060102")
	}
	perPage := 100

	accountId := os.Getenv("ACCOUNT_ID")
	accessToken := os.Getenv("TOKEN")

	ctx := context.Background()

	log.Debug("Fetching Harvest entries via HTTP API...")
	nextPage := 1
	entries := []*harvest.TimeEntry{}

	client := harvest.NewClient(accessToken, accountId)

	for {
		log.WithFields(log.Fields{
			"page": nextPage,
		}).Trace("Querying page...")

		opts := &harvest.TimeEntriesListOptions{
			From:        from,
			To:          to,
			ListOptions: harvest.ListOptions{Page: nextPage, PerPage: perPage},
		}

		timeEntries, r, err := client.TimeEntries.List(ctx, opts)
		if err != nil {
			return entries, err
		}

		entries = append(entries, timeEntries...)

		if r.NextPage == nextPage {
			break
		}
		nextPage = r.NextPage
	}
	log.Debug("Fetched Harvest entries.")

	return entries, nil
}

func parseOutputFormat(of string) (OutputFormat, error) {
	switch strings.ToLower(of) {
	case "csv":
		return CSVOutputFormat, nil
	case "text":
		return TextOutputFormat, nil
	}

	var o OutputFormat
	return o, fmt.Errorf("not a valid output format: %q", of)
}

const (
	TextOutputFormat OutputFormat = iota
	CSVOutputFormat
)
