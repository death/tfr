package main

import (
	"fmt"

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

	stats, err := store.Statistics()
	if err != nil {
		return err
	}

	fmt.Printf("%5d Unread\n", stats.Unread)
	fmt.Printf("%5d In progress\n", stats.InProgress)
	fmt.Printf("%5d Read\n", stats.Read)

	return nil
}
