package main

import "github.com/death/tfr/db"

type IndexCommand struct {
}

var index IndexCommand

func init() {
	parser.AddCommand("index",
		"Index text files",
		"Index text files.",
		&index)
}

func (c *IndexCommand) Execute(args []string) error {
	return db.Build(options.BaseDir, options.DBFile)
}
