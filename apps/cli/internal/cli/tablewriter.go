package cli

import (
	"fmt"
	"io"
	"strings"
)

// TableWriter writes aligned table output where ANSI escape codes don't affect alignment.
// It accepts both plain-text and colored versions of each cell, computing column widths
// from the plain-text values and emitting the colored text with correct padding.
type TableWriter struct {
	rows []tableRow
	gap  string
}

type tableRow struct {
	plain   []string
	colored []string
	isSep   bool
}

// NewTableWriter creates a TableWriter with the default "  " column gap.
func NewTableWriter() *TableWriter {
	return &TableWriter{gap: "  "}
}

// AddHeader adds a plain-text header row (no coloring).
func (tw *TableWriter) AddHeader(cols []string) {
	tw.rows = append(tw.rows, tableRow{plain: cols, colored: cols})
}

// AddSeparator adds a row of dashes sized to match column widths.
// The actual dash strings are computed in Flush.
func (tw *TableWriter) AddSeparator() {
	tw.rows = append(tw.rows, tableRow{isSep: true})
}

// AddRow adds a data row. plain holds the visible text for width calculation;
// colored holds the text to render (may contain ANSI codes).
func (tw *TableWriter) AddRow(plain, colored []string) {
	tw.rows = append(tw.rows, tableRow{plain: plain, colored: colored})
}

// Flush computes column widths and writes all rows to w.
func (tw *TableWriter) Flush(w io.Writer) {
	if len(tw.rows) == 0 {
		return
	}

	widths := tw.columnWidths()

	for _, row := range tw.rows {
		if row.isSep {
			tw.writeSeparator(w, widths)
			continue
		}
		tw.writeRow(w, row, widths)
	}
}

// columnWidths returns the max visible width for each column.
func (tw *TableWriter) columnWidths() []int {
	var maxCols int
	for _, row := range tw.rows {
		if !row.isSep && len(row.plain) > maxCols {
			maxCols = len(row.plain)
		}
	}

	widths := make([]int, maxCols)
	for _, row := range tw.rows {
		if row.isSep {
			continue
		}
		for i, cell := range row.plain {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	return widths
}

func (tw *TableWriter) writeSeparator(w io.Writer, widths []int) {
	cells := make([]string, len(widths))
	for i, width := range widths {
		cells[i] = strings.Repeat("-", width)
	}
	fmt.Fprintln(w, strings.Join(cells, tw.gap))
}

func (tw *TableWriter) writeRow(w io.Writer, row tableRow, widths []int) {
	cells := make([]string, len(widths))
	for i := range widths {
		var plain, colored string
		if i < len(row.plain) {
			plain = row.plain[i]
		}
		if i < len(row.colored) {
			colored = row.colored[i]
		}
		padding := widths[i] - len(plain)
		if padding < 0 {
			padding = 0
		}
		cells[i] = colored + strings.Repeat(" ", padding)
	}
	fmt.Fprintln(w, strings.Join(cells, tw.gap))
}

// StripANSI removes ANSI escape sequences from a string.
func StripANSI(s string) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && !((s[j] >= 'A' && s[j] <= 'Z') || (s[j] >= 'a' && s[j] <= 'z')) {
				j++
			}
			if j < len(s) {
				j++ // skip the terminator
			}
			i = j
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}
