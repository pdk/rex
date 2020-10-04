package rex

// JSON readers/writers

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

func ReadJSONLines(input io.Reader) chan Record {

	c := make(chan Record)

	go func() {
		defer close(c)

		scanner := bufio.NewScanner(input)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Text()
			var obj interface{}
			err := json.Unmarshal([]byte(line), &obj)
			if err != nil {
				log.Fatalf("ReadJSONLines failed to parse JSON: %v", err)
			}

			c <- NewRecord(obj)
		}

		err := scanner.Err()
		if err != nil {
			log.Fatalf("ReadJSONLines failed reading input: %v", err)
		}
	}()

	return c
}

// WriteJSONLinesProcessor writes each object as JSON on a line by itself.
func WriteJSONLinesProcessor(output io.Writer) RecordProcessor {

	return func(r Record) Record {

		result, err := json.Marshal(r.values)
		if err != nil {
			log.Fatalf("failed to output JSON for object: %v", err)
		}

		_, err = output.Write(result)
		if err != nil {
			log.Fatalf("failed to write output: %v", err)
		}

		_, err = output.Write([]byte("\n"))
		if err != nil {
			log.Fatalf("failed to write output: %v", err)
		}

		// pass it along
		return r
	}
}

// WriteJSONLines outputs one JSON object per line.
func WriteJSONLines(output io.Writer) RecordCollector {

	writer := WriteJSONLinesProcessor(output)

	return func(c chan Record) {
		for r := range c {
			writer(r)
		}
	}
}
