package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	BaseDir string `long:"base-dir" description:"Base directory" default:"/media/1984/Documents/textfiles/"`
	DBFile  string `long:"db-file" description:"Database file name" default:"/media/1984/Documents/textfiles/tfr.db"`
}

var (
	options Options
	parser  = flags.NewParser(&options, flags.Default)
)

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
