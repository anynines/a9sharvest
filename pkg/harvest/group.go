package harvest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

type Content struct {
	TimeEntries []struct {
		Hours float64
		Notes string
	} `json:"time_entries"`
}

func Group(verboseFlag bool) error {
	if verboseFlag {
		log.SetLevel(log.DebugLevel)
	}

	err := CheckEnvVariables()
	if err != nil {
		return err
	}

	log.Debug("Querying harvest...")
	data, err := fetchData()
	if err != nil {
		return err
	}
	log.Debug("Queried harvest.")

	var content Content
	err = json.Unmarshal(data, &content)
	if err != nil {
		return err
	}

	for _, v := range content.TimeEntries {
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

// FIXME iterate over multiple pages -> pagination
// FIXME hardcoded from, to values
func fetchData() ([]byte, error) {
	v := url.Values{}
	v.Set("from", "20200901")
	v.Set("to", "20200911")
	v.Set("page", "1")

	// https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/#list-all-time-entries
	url := "https://api.harvestapp.com/v2/time_entries" + "?" //+ v.Encode()
	body, err := HttpGet(url, os.Getenv("ACCOUNT_ID"), os.Getenv("TOKEN"))
	if err != nil {
		return body, err
	}

	return body, nil
}
