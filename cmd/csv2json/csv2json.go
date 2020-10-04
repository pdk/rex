package main

import (
	"os"

	"github.com/pdk/rex"
)

func main() {
	rex.Unue(rex.ReadCSV(os.Stdin)).
		Kolekti(rex.WriteJSON(os.Stdout))
}
