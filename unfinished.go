package main

import "github.com/death/tfr/db"

type UnfinishedCommand struct {
}

var unfinished UnfinishedCommand

func init() {
	parser.AddCommand("unfinished",
		"Mark latest article read as unfinished",
		"Mark latest article read as unfinished.",
		&unfinished)
}

func (c *UnfinishedCommand) Execute(args []string) error {
	store, err := db.NewStore(options.DBFile)
	if err != nil {
		return err
	}
	defer store.Close()

	article, err := store.LatestReadArticle()
	if err != nil {
		return err
	}

	err = store.MarkAsUnfinished(article.ID)
	if err != nil {
		return err
	}

	return nil
}
