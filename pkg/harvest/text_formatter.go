package harvest

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type TextFormatter struct {
	Result *Result
}

func NewTextFormatter(result *Result) *TextFormatter {
	return &TextFormatter{
		Result: result,
	}
}

func (f *TextFormatter) Output() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)

	for k, v := range f.Result.GroupedByTag {
		table.Append([]string{k, fmt.Sprintf("%.2f", v)})
	}

	table.Render()
}
