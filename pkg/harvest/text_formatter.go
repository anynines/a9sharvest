package harvest

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type TextFormatter struct {
	Stats *Stats
}

func NewTextFormatter(stats *Stats) *TextFormatter {
	return &TextFormatter{
		Stats: stats,
	}
}

func (f *TextFormatter) Output() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)

	for k, v := range f.Stats.GroupedByTag {
		table.Append([]string{k, fmt.Sprintf("%.2f", v)})
	}

	table.Render()
}
