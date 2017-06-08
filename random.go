package main

import (
	"fmt"
	"log"

	"github.com/death/tfr/db"
)

type RandomCommand struct {
}

var random RandomCommand

func init() {
	parser.AddCommand("random",
		"Read a random unread text file",
		"Read a random unread text file.",
		&random)
}

func (c *RandomCommand) Execute(args []string) error {
	store, err := db.NewStore(options.DBFile)
	if err != nil {
		return err
	}
	defer store.Close()

	article, err := store.RandomArticle()
	if err != nil {
		return err
	}

	section, err := store.FindSection(article.SectionID)
	if err != nil {
		return err
	}

	fmt.Printf("[ %s : %s ]\n", section.Path, article.Path)

	cat.ArchiveFile = section.Path
	cat.TextfilePath = article.Path
	err = cat.Execute(nil)
	if err != nil {
		return err
	}

	err = store.MarkAsRead(article.ID)
	if err != nil {
		log.Printf("Could not mark article %d as read: %v", article.ID, err)
	}

	return nil
}
