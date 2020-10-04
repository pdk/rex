package rex

import (
	"encoding/csv"
	"io"
	"log"
)

// ReadCSV is a RecordProducer for reading CSV inputs.
func ReadCSV(input io.Reader) chan Record {

	c := make(chan Record)

	reader := csv.NewReader(input)

	go func() {
		defer close(c)

		var headers []string

		for i := 0; ; i++ {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("failed to read csv: %s", err)
			}

			if i == 0 {
				// first line. gather headers
				headers = record
			} else {
				vals := map[string]interface{}{}
				for j := 0; j < len(headers); j++ {
					vals[headers[j]] = record[j]
				}

				c <- Record{
					values: vals,
				}
			}
		}
	}()

	return c
}

func WriteCSV(output io.Writer) RecordCollector {

	writer := csv.NewWriter(output)

	return func(c chan Record) {
		defer writer.Flush()

		first, ok := <-c
		if !ok {
			return
		}

		headers := []string{}
		for k := range first.values {
			headers = append(headers, k)
		}

		err := writer.Write(headers)
		if err != nil {
			log.Fatalf("can't write headers: %v", err)
		}

		writeRecordCSV(writer, headers, first)

		for r := range c {
			writeRecordCSV(writer, headers, r)
		}
	}
}

func writeRecordCSV(writer *csv.Writer, headers []string, r Record) {

	row := []string{}
	for _, h := range headers {
		row = append(row, r.String(h))
	}

	err := writer.Write(row)
	if err != nil {
		log.Fatalf("can't write row: %v", err)
	}
}
