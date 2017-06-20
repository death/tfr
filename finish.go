package main

import (
	"errors"

	"github.com/death/tfr/db"
)

type FinishCommand struct {
}

var finish FinishCommand

func init() {
	parser.AddCommand("finish",
		"Mark latest viewed article as finished",
		`Remove the article from the unfinished queue and from
participation in random picks.`,
		&finish)
}

func (c *FinishCommand) Execute(args []string) error {
	store, err := db.NewStore(options.DBFile)
	if err != nil {
		return err
	}
	defer store.Close()

	articleID, err := store.LatestViewedArticleID()
	if err != nil {
		return err
	}

	if articleID == db.NoArticle {
		return errors.New("no article was viewed")
	}

	err = store.MarkAsRead(articleID)
	if err != nil {
		return err
	}

	return nil
}
