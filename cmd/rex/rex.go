package main

import (
	"log"
	"os"

	"github.com/pdk/rex"
)

func main() {

	if len(os.Args) != 3 {
		log.Printf("usage: rex in out")
		log.Fatalf("in/out: {file.}csv or {file.}json or {file.}jsonl")
	}

	rex.Unue(rex.ProducerFrom(os.Args[1])).
		// Ekskludi(func(r rex.Record) bool {
		// 	return r.Int("beta") == 10
		// }).
		// Kaj(func(r rex.Record) rex.Record {
		// 	r.Set("foo", r.Int("beta")*100)
		// 	return r
		// }).
		Kolekti(rex.OutputTo(os.Args[2]))
}
