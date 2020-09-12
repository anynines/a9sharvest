package harvest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type TimeEntry struct {
	Hours float64
	Notes string
}

type Content struct {
	TimeEntries []TimeEntry `json:"time_entries"`
	NextPage    *int        `json:"next_page"`
}

func Group(verboseFlag bool) error {
	if verboseFlag {
		log.SetLevel(log.DebugLevel)
	}

	err := CheckEnvVariables()
	if err != nil {
		return err
	}

	entries, err := fetchTimeEntries()
	if err != nil {
		return err
	}

	for _, v := range entries {
		fmt.Printf("%v : %v\n", v.Hours, v.Notes)
	}

	log.Debug("Done.")

	return nil
}

func CheckEnvVariables() error {
	keys := []string{"ACCOUNT_ID", "TOKEN"}

	for _, key := range keys {
		if len(os.Getenv(key)) < 1 {
			return fmt.Errorf("You MUST set the environment variable %s", key)
		}
	}

	return nil
}

func fetchTimeEntries() ([]TimeEntry, error) {
	entries := []TimeEntry{}

	nextPage := 1
	for {
		log.WithFields(log.Fields{
			"page": nextPage,
		}).Debug("Querying page...")
		data, err := fetchData(nextPage)
		if err != nil {
			return entries, err
		}
		log.Debug("Queried page.")

		var content Content
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

	return entries, nil
}

// FIXME hardcoded from, to values
func fetchData(page int) ([]byte, error) {
	v := url.Values{}
	v.Set("from", "20200901")
	v.Set("to", "20200911")
	v.Set("page", strconv.Itoa(page))
	v.Set("per_page", "20")

	// https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/#list-all-time-entries
	url := "https://api.harvestapp.com/v2/time_entries" + "?" + v.Encode()
	return HttpGet(url, os.Getenv("ACCOUNT_ID"), os.Getenv("TOKEN"))
}
