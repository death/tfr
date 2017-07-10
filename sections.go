package main

import (
	"fmt"

	"github.com/death/tfr/db"
)

type SectionsCommand struct {
}

var sections SectionsCommand

func init() {
	parser.AddCommand("sections",
		"List sections",
		"List sections.",
		&sections)
}

func (c *SectionsCommand) Execute(args []string) error {
	store, err := db.NewStore(options.DBFile)
	if err != nil {
		return err
	}
	defer store.Close()

	sections, err := store.ListSections()
	if err != nil {
		return err
	}

	fmt.Println(db.AnySectionLabel)
	for _, section := range sections {
		fmt.Println(section.Label)
	}

	return nil
}
