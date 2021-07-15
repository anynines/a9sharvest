package harvest

import (
	"os"
	"strconv"
	"strings"

	"github.com/gogolok/go-harvest/harvest"
	log "github.com/sirupsen/logrus"
)

type ReportInterface interface {
	Output()
}

type Report struct {
	TimeEntries  []*harvest.TimeEntry
	matcher      EntryMatcher
	groupedByTag map[string]float64
	skipUnknown  bool
}

func NewReport(timeEntries []*harvest.TimeEntry, matcher EntryMatcher, skipUnknown bool) Report {
	return Report{
		TimeEntries: timeEntries,
		matcher:     matcher,
		skipUnknown: skipUnknown,
	}
}

func (r *Report) Run() {
	skip_project_ids := strings.Split(os.Getenv("SKIP_PROJECT_IDS"), ",")
	skip_project_ids_map := map[string]int{}
	for _, v := range skip_project_ids {
		skip_project_ids_map[v] = 1
	}

	allowed_project_ids_enabled := len(os.Getenv("ALLOWED_PROJECT_IDS")) > 0
	allowed_project_ids := strings.Split(strings.TrimSpace(os.Getenv("ALLOWED_PROJECT_IDS")), ",")
	allowed_project_ids_map := map[string]int{}
	for _, v := range allowed_project_ids {
		allowed_project_ids_map[v] = 1
	}

	allowed_user_ids_enabled := len(os.Getenv("ALLOWED_USER_IDS")) > 0
	allowed_user_ids := strings.Split(os.Getenv("ALLOWED_USER_IDS"), ",")
	allowed_user_ids_map := map[string]int{}
	for _, v := range allowed_user_ids {
		allowed_user_ids_map[v] = 1
	}

	allowed_task_names_enabled := len(os.Getenv("ALLOWED_TASK_NAMES")) > 0
	allowed_task_names := strings.Split(os.Getenv("ALLOWED_TASK_NAMES"), ",")
	allowed_task_names_map := map[string]int{}
	for _, v := range allowed_task_names {
		allowed_task_names_map[v] = 1
	}

	TAG_UNKNOWN := "[unknown]"

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
			"task-id":      v.Task.Id,
			"task-name":    v.Task.Name,
		}
		log.WithFields(logFields).Trace("time entry")

		if _, ok := skip_project_ids_map[strconv.Itoa(v.Project.Id)]; ok {
			log.WithFields(logFields).Trace("Skipped because of project id")
			continue
		}

		if allowed_project_ids_enabled {
			if _, ok := allowed_project_ids_map[strconv.Itoa(v.Project.Id)]; !ok {
				log.WithFields(logFields).Trace("Skipped because of project id (not allowed)")
				continue
			}
		}

		if allowed_user_ids_enabled {
			if _, ok := allowed_user_ids_map[strconv.Itoa(v.User.Id)]; !ok {
				log.WithFields(logFields).Trace("Skipped because of user id")
				continue
			}
		}

		if allowed_task_names_enabled {
			if _, ok := allowed_task_names_map[v.Task.Name]; !ok {
				log.WithFields(logFields).Trace("Skipped because of task name")
				continue
			}
		}

		matched, tag := r.matcher.Match(v.Notes)
		if matched {
			grouped_by_tag[tag] += v.Hours
			log.WithFields(logFields).Info("time entry matched")
		} else {
			log.WithFields(logFields).Debug("New [unknown] entry")
			if !r.skipUnknown {
				grouped_by_tag[TAG_UNKNOWN] += v.Hours
			}
		}
	}

	r.groupedByTag = grouped_by_tag
}

func (r *Report) Stats() *Stats {
	return NewStats(r.groupedByTag)
}
