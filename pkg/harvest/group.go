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

type Project struct {
	Id   int
	Name string
}

type TimeEntry struct {
	Id      int
	Hours   float64
	Notes   string
	Project Project
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

	skip_project_ids := strings.Split(os.Getenv("SKIP_PROJECT_IDS"), ",")
	skip_project_ids_map := map[string]int{}
	for _, v := range skip_project_ids {
		skip_project_ids_map[v] = 1
	}

	TAG_UNKNOWN := "[unknown]"
	tags := strings.Split(os.Getenv("TAGS"), ",")
	log.WithFields(log.Fields{
		"tags": tags,
	}).Debug("Set up tags.")
	grouped_by_tags := make(map[string]float64)

	for _, v := range entries {
		log.WithFields(log.Fields{
			"id":           v.Id,
			"project-id":   v.Project.Id,
			"project-name": v.Project.Name,
			"hours":        v.Hours,
			"notes":        v.Notes,
		}).Debug("time entry")

		if _, ok := skip_project_ids_map[strconv.Itoa(v.Project.Id)]; ok {
			log.WithFields(log.Fields{
				"project-id": v.Project.Id,
			}).Debug("Skipped because of project id")
			continue
		}

		matched := false

		for _, t := range tags {
			if strings.Contains(v.Notes, t) {
				grouped_by_tags[t] += v.Hours
				matched = true
				break
			}
		}

		if !matched {
			log.WithFields(log.Fields{
				"Notes": v.Notes,
			}).Debug("New [unknown] entry")
			grouped_by_tags[TAG_UNKNOWN] += v.Hours
		}
	}

	for k, v := range grouped_by_tags {
		fmt.Printf("%v = %v\n", k, v)
	}

	log.Debug("Done.")

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
	v.Set("per_page", "20")

	// https://help.getharvest.com/api-v2/timesheets-api/timesheets/time-entries/#list-all-time-entries
	url := "https://api.harvestapp.com/v2/time_entries" + "?" + v.Encode()
	return HttpGet(url, os.Getenv("ACCOUNT_ID"), os.Getenv("TOKEN"))
}
