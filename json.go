package rex

// JSON readers/writers

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

// ReadJSON will read a JSON document. JSON should contain an array of objects,
// otherwise, program will terminate.
func ReadJSON(input io.Reader) chan Record {

	c := make(chan Record)

	go func() {
		defer close(c)

		dec := json.NewDecoder(input)

		t, err := dec.Token()
		if err != nil {
			log.Fatalf("failed to read stdin: %v", err)
		}

		if t != json.Delim('[') {
			log.Fatalf("expected to read JSON array starting with \"[\", but found %v", t)
		}

		for dec.More() {

			var obj interface{}
			err := dec.Decode(&obj)
			if err != nil {
				log.Fatalf("failed to ready JSON object: %v", err)
			}

			c <- NewRecord(obj)
		}

		_, err = dec.Token()
		if err != nil {
			log.Fatalf("failed to read final \"]\": %v", err)
		}
	}()

	return c
}

// WriteJSON writes an array of records.
func WriteJSON(output io.Writer) RecordCollector {

	writer := bufio.NewWriter(output)

	return func(c chan Record) {
		defer writer.Flush()

		_, err := writer.WriteString("[\n    ")
		if err != nil {
			log.Fatalf("WriteJSON failed to write start of output: %v", err)
		}

		i := 0
		for r := range c {
			if i > 0 {
				_, err := writer.WriteString(",\n    ")
				if err != nil {
					log.Fatalf("WriteJSON failed to write object separator: %v", err)
				}
			}
			i++

			result, err := json.MarshalIndent(r.values, "    ", "    ")
			if err != nil {
				log.Fatalf("can't unmarshal data: %s", err)
			}

			_, err = writer.Write(result)
			if err != nil {
				log.Fatalf("WriteJSON failed to write start of output: %v", err)
			}
		}

		_, err = writer.WriteString("\n]\n")
		if err != nil {
			log.Fatalf("WriteJSON failed to write end of output: %v", err)
		}
	}
}
