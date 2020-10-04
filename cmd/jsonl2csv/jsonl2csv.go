package main

import (
	"os"

	"github.com/pdk/rex"
)

func main() {
	rex.Unue(rex.ReadJSONLines(os.Stdin)).
		Kolekti(rex.WriteCSV(os.Stdout))
}
