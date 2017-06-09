package main

import (
	"fmt"
	"log"

	"github.com/death/tfr/db"
)

type ResumeCommand struct {
}

var resume ResumeCommand

func init() {
	parser.AddCommand("resume",
		"Resume reading oldest unfinished article",
		"Resume reading oldest unfinished article.",
		&resume)
}

func (c *ResumeCommand) Execute(args []string) error {
	store, err := db.NewStore(options.DBFile)
	if err != nil {
		return err
	}
	defer store.Close()

	article, err := store.OldestUnfinishedArticle()
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
