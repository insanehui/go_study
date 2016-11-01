package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// flags are options, with pre "-"
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	// args are arguments without "-"
	name = kingpin.Arg("name", "Name of user.").Required().String()
)

func main() {
	kingpin.Parse()
	fmt.Printf("%v, %s\n", *verbose, *name)
}
