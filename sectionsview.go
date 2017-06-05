package main

import (
	"fmt"

	"github.com/death/tfr/db"
	"github.com/nsf/termbox-go"
)

type SectionsView struct {
	list *ListView
}

func NewSectionsView(store *db.Store) *SectionsView {
	sections := store.AllSections()
	items := make([]string, len(sections))
	for i, section := range sections {
		items[i] = section.Label
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

	return &SectionsView{
		list: list,
	}
}

func (v *SectionsView) Draw() {
	v.list.Draw()
	v.drawStatusLine()
}

func (v *SectionsView) drawStatusLine() {
	w, h := termbox.Size()
	x, y := 0, h-1

	fg := termbox.ColorDefault
	bg := termbox.ColorDefault

	status := fmt.Sprintf("This is the status line.. selected item is %d/%d", v.list.SelectedItem()+1, v.list.NumItems())

	for _, c := range status {
		if x >= w {
			break
		}
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
	for x < w {
		termbox.SetCell(x, y, ' ', fg, bg)
		x++
	}
}

func (v *SectionsView) HandleEvent(ev termbox.Event) bool {
	return v.list.HandleEvent(ev)
}
