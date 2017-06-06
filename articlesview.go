package main

import (
	"github.com/death/tfr/db"
	"github.com/nsf/termbox-go"
)

type ArticlesView struct {
	list     *ListView
	store    *db.Store
	articles *[]db.Article
}

func NewArticlesView(store *db.Store, sectionID int) *ArticlesView {
	articles := store.ArticlesForSection(sectionID)
	items := make([]string, len(articles))
	for i, article := range articles {
		items[i] = article.Label
	}

	w, h := termbox.Size()

	list := &ListView{
		StartX: 0,
		EndX:   w,
		StartY: 0,
		EndY:   h - 1,

		ForegroundColor: termbox.ColorDefault,
		BackgroundColor: termbox.ColorDefault,

		SelectForegroundColor: termbox.ColorWhite,
		SelectBackgroundColor: termbox.ColorBlue,

		Items: items,
	}

	return &ArticlesView{
		list:     list,
		store:    store,
		articles: &articles,
	}
}

func (v *ArticlesView) Draw() {
	v.list.Draw()
}

func (v *ArticlesView) HandleEvent(ev termbox.Event) (bool, View) {
	if ev.Type == termbox.EventKey {
		if ev.Key == termbox.KeyEnter {
			// Show article
		}
	}
	return v.list.HandleEvent(ev)
}
