package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type SectionsView struct {
	list *ListView
}

func NewSectionsView() *SectionsView {
	sections := []string{
		"100",
		"Adventure",
		"Anarchy",
		"Apple",
		"Art",
		"Bbs",
		"Computers",
		"Conspiracy",
		"Drugs",
		"Food",
		"Fun",
		"Games",
		"Groups",
		"Hacking",
		"Hamradio",
		"Holiday",
		"Humor",
		"Internet",
		"Law",
		"Magazines",
		"Media",
		"Messages",
		"Music",
		"News",
		"Occult",
		"Phreak",
		"Piracy",
		"Politics",
		"Programming",
		"Reports",
		"Rpg",
		"Science",
		"Sex",
		"Sf",
		"Stories",
		"Survival",
		"Ufo",
		"Uploads",
		"Virus",
	}

	w, _ := termbox.Size()

	list := &ListView{
		StartX: 0,
		EndX:   w,
		StartY: 0,
		EndY:   10,

		ForegroundColor: termbox.ColorDefault,
		BackgroundColor: termbox.ColorDefault,

		SelectForegroundColor: termbox.ColorWhite,
		SelectBackgroundColor: termbox.ColorBlue,

		Items: sections,
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
