package rex

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// MaxColWidth limits how wide a column can be.
const MaxColWidth = 80

// WriteTabular writes a stream of Records in plaintext table format.
func WriteTabular(output io.Writer) RecordCollector {

	return func(c chan Record) {

		var records []Record
		columns := map[string]bool{}
		widths := map[string]int{}

		for r := range c {
			records = append(records, r)

			for c := range r.values {
				columns[c] = true
				widths[c] = max(widths[c], len(r.String(c)))
			}
		}

		cols := []string{}
		for c := range columns {
			cols = append(cols, c)
			widths[c] = max(widths[c], len(c))
		}

		writeTable(output, records, cols, widths)
	}
}

// writeTable prints a set of records as a plaintext table.
func writeTable(output io.Writer, recs []Record, columns []string, widths map[string]int) {

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	for _, col := range columns {
		fmt.Printf("%-*s  ", min(MaxColWidth, widths[col]), col)
	}
	fmt.Println()

	for _, col := range columns {
		fmt.Printf("%s  ", strings.Repeat("-", min(MaxColWidth, widths[col])))
	}
	fmt.Println()

	for _, row := range recs {
		for _, col := range columns {
			fmt.Printf("%-*s  ", min(MaxColWidth, widths[col]), trunc(row.String(col), MaxColWidth))
		}
		fmt.Println()
	}
}

// trunc returns the leading n characters of a string, or the whole string if shorter than n.
func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[0:n]
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
