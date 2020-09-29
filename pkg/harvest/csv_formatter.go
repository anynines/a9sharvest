package harvest

import (
	"encoding/csv"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type CSVFormatter struct {
	Stats *Stats
}

func NewCSVFormatter(stats *Stats) *CSVFormatter {
	return &CSVFormatter{
		Stats: stats,
	}
}

func (f *CSVFormatter) Output() {
	w := csv.NewWriter(os.Stdout)

	record := []string{"Tag", "Hours", "Percentage"}
	if err := w.Write(record); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for tag, v := range f.Stats.GroupedByTag {
		p := f.Stats.PercentageForTag(tag)
		record := []string{tag, fmt.Sprintf("%.2f", v), fmt.Sprintf("%.2f", p)}
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()
}
