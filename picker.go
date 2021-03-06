package rex

import (
	"log"
	"os"
	"path/filepath"
)

// ProducerFrom selects the type of input (and source) from the path/extension.
func ProducerFrom(path string) chan Record {

	switch path {
	case "csv":
		return ReadCSV(os.Stdin)
	case "json":
		return ReadJSON(os.Stdin)
	case "jsonl":
		return ReadJSONLines(os.Stdin)
	}

	in, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open input file %s: %v", path, err)
	}

	switch filepath.Ext(path) {
	case ".csv":
		return ReadCSV(in)
	case ".json":
		return ReadJSON(in)
	case ".jsonl":
		return ReadJSONLines(in)
	}

	log.Fatalf("unrecognized input file extension: %s", filepath.Ext(path))
	return nil
}

// OutputTo selects type of output handling and destination based on path/extensions.
func OutputTo(path string) func(chan Record) {

	switch path {
	case "csv":
		return WriteCSV(os.Stdout)
	case "json":
		return WriteJSON(os.Stdout)
	case "jsonl":
		return WriteJSONLines(os.Stdout)
	case "tab", "tabular", "table":
		return WriteTabular(os.Stdout)
	}

	_, err := os.Stat(path)
	if err == nil {
		log.Fatalf("output file %s already exists - refusing to overwrite", path)
	}

	out, err := os.Create(path)
	if err != nil {
		log.Fatalf("failed to create output file %s: %v", path, err)
	}

	switch filepath.Ext(path) {
	case ".csv":
		return WriteCSV(out)
	case ".json":
		return WriteJSON(out)
	case ".jsonl":
		return WriteJSONLines(out)
	case ".tab", ".tabular", ".table":
		return WriteTabular(out)
	}

	log.Fatalf("unrecognized output file extension: %s", filepath.Ext(path))
	return nil
}
