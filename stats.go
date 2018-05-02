package main

import (
	"fmt"
	"sort"

	"github.com/death/tfr/db"
)

type StatsCommand struct {
}

var stats StatsCommand

func init() {
	parser.AddCommand("stats",
		"Show statistics",
		"Show statistics.",
		&stats)
}

func (c *StatsCommand) Execute(args []string) error {
	store, err := db.NewStore(options.DBFile)
	if err != nil {
		return err
	}
	defer store.Close()

	statsMap, err := store.Statistics()
	if err != nil {
		return err
	}

	var sections []string
	for section := range statsMap {
		sections = append(sections, section)
	}
	sort.Strings(sections)

	fmt.Printf("%-20s %6s %6s %6s\n", "SECTION", "UNREAD", "INPROG", "READ")

	for _, section := range sections {
		stats := statsMap[section]
		fmt.Printf("%-20s %6d %6d %6d\n", section, stats.Unread, stats.InProgress, stats.Read)
	}

	return nil
}
