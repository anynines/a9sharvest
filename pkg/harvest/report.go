package harvest

import (
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Report struct {
	TimeEntries  []TimeEntry
	groupedByTag map[string]float64
}

func NewReport(timeEntries []TimeEntry) Report {
	return Report{
		TimeEntries: timeEntries,
	}
}

func (r *Report) Run() {
	skip_project_ids := strings.Split(os.Getenv("SKIP_PROJECT_IDS"), ",")
	skip_project_ids_map := map[string]int{}
	for _, v := range skip_project_ids {
		skip_project_ids_map[v] = 1
	}

	allowed_user_ids_enabled := len(os.Getenv("ALLOWED_USER_IDS")) > 0
	allowed_user_ids := strings.Split(os.Getenv("ALLOWED_USER_IDS"), ",")
	allowed_user_ids_map := map[string]int{}
	for _, v := range allowed_user_ids {
		allowed_user_ids_map[v] = 1
	}

	TAG_UNKNOWN := "[unknown]"
	tags := strings.Split(os.Getenv("TAGS"), ",")
	log.WithFields(log.Fields{
		"tags": tags,
	}).Debug("Set up tags.")
	grouped_by_tag := make(map[string]float64)

	for _, v := range r.TimeEntries {
		logFields := log.Fields{
			"id":           v.Id,
			"spent_date":   v.SpentDate,
			"hours":        v.Hours,
			"notes":        v.Notes,
			"project-id":   v.Project.Id,
			"project-name": v.Project.Name,
			"user-id":      v.User.Id,
			"user-name":    v.User.Name,
		}
		log.WithFields(logFields).Trace("time entry")

		if _, ok := skip_project_ids_map[strconv.Itoa(v.Project.Id)]; ok {
			log.WithFields(logFields).Trace("Skipped because of project id")
			continue
		}

		if allowed_user_ids_enabled {
			if _, ok := allowed_user_ids_map[strconv.Itoa(v.User.Id)]; !ok {
				log.WithFields(logFields).Trace("Skipped because of user id")
				continue
			}
		}

		matched := false

		for _, t := range tags {
			if strings.Contains(v.Notes, t) {
				grouped_by_tag[t] += v.Hours
				matched = true
				break
			}
		}

		if !matched {
			log.WithFields(logFields).Debug("New [unknown] entry")
			grouped_by_tag[TAG_UNKNOWN] += v.Hours
		}
	}

	r.groupedByTag = grouped_by_tag
}

func (r *Report) Stats() *Stats {
	return NewStats(r.groupedByTag)
}