package harvest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type OutputFormat uint32

type Project struct {
	Id   int
	Name string
}

type User struct {
	Id   int
	Name string
}

type TimeEntry struct {
	Id        int
	Hours     float64
	Notes     string
	Project   Project
	User      User
	SpentDate string `json:"spent_date"`
}

type TimeEntriesPage struct {
	TimeEntries []TimeEntry `json:"time_entries"`
	NextPage    *int        `json:"next_page"`
}

type Result struct {
	GroupedByTag map[string]float64
}

func Group(verboseFlag bool, outputFlag string) error {
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

	err := CheckEnvVariables()
	if err != nil {
		return err
	}

	outputFormat, err := parseOutputFormat(outputFlag)
	if err != nil {
		log.Fatalln(err)
	}

	entries, err := fetchTimeEntries()
	if err != nil {
		return err
	}

	report := NewReport(entries)
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

func CheckEnvVariables() error {
	keys := []string{"ACCOUNT_ID", "TOKEN", "TAGS"}

	for _, key := range keys {
		if len(os.Getenv(key)) < 1 {
			return fmt.Errorf("You MUST set the environment variable %s", key)
		}
	}

	return nil
}

func fetchTimeEntries() ([]TimeEntry, error) {
	entries := []TimeEntry{}

	log.Debug("Fetching Harvest entries via HTTP API...")
	nextPage := 1
	for {
		log.WithFields(log.Fields{
			"page": nextPage,
		}).Trace("Querying page...")
		data, err := fetchData(nextPage)
		if err != nil {
			return entries, err
		}

		var content TimeEntriesPage
		err = json.Unmarshal(data, &content)
		if err != nil {
			return entries, err
		}
		entries = append(entries, content.TimeEntries...)

		if content.NextPage == nil {
			break
		}
		nextPage = *content.NextPage
	}
	log.Debug("Fetched Harvest entries.")

	return entries, nil
}

func fetchData(page int) ([]byte, error) {
	to := os.Getenv("TO")
	if len(to) < 1 {
		to = time.Now().Format("20060102")
	}

	from := os.Getenv("FROM")
	if len(from) < 1 {
		from = time.Now().AddDate(0, 0, -14).Format("20060102")
	}

	v := url.Values{}
	v.Set("from", from)
	v.Set("to", to)
	v.Set("page", strconv.Itoa(page))
	v.Set("per_page", "100")

	// https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/#list-all-time-entries
	url := "https://api.harvestapp.com/v2/time_entries" + "?" + v.Encode()
	return HttpGet(url, os.Getenv("ACCOUNT_ID"), os.Getenv("TOKEN"))
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
