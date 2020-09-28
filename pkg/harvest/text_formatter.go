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

	table.SetHeader([]string{"Tag", "Hours", "%"})
	for tag, v := range f.Stats.GroupedByTag {
		p := f.Stats.PercentageForTag(tag)
		table.Append([]string{tag, fmt.Sprintf("%.2f", v), fmt.Sprintf("%.2f", p)})
	}

	table.Render()
}
