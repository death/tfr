package main

import (
	"errors"
	"fmt"

	"github.com/death/tfr/db"
)

const (
	PickRandom = iota
	PickUnfinished
)

func ViewArticle(which int) error {
	store, err := db.NewStore(options.DBFile)
	if err != nil {
		return err
	}
	defer store.Close()

	var article *db.Article
	switch which {
	case PickRandom:
		article, err = store.RandomArticle()
	case PickUnfinished:
		article, err = store.OldestUnfinishedArticle()
	default:
		return errors.New("weird article picker")
	}
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

	err = store.MarkAsUnfinished(article.ID)
	if err != nil {
		return err
	}

	err = store.SetLatestViewedArticleID(article.ID)
	if err != nil {
		return err
	}

	return nil
}
