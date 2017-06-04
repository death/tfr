package main

import "github.com/nsf/termbox-go"

type ListView struct {
	StartX int
	EndX   int
	StartY int
	EndY   int

	ForegroundColor termbox.Attribute
	BackgroundColor termbox.Attribute

	SelectForegroundColor termbox.Attribute
	SelectBackgroundColor termbox.Attribute

	Items []string

	top      int
	selected int
}

func (v *ListView) Draw() {
	if len(v.Items) == 0 {
		return
	}

	y := v.StartY

	for i, item := range v.Items[v.top:] {
		if y >= v.EndY {
			break
		}

		k := i + v.top

		var fg, bg termbox.Attribute
		if k == v.selected {
			fg = v.SelectForegroundColor
			bg = v.SelectBackgroundColor
		} else {
			fg = v.ForegroundColor
			bg = v.BackgroundColor
		}

		v.drawItem(y, fg, bg, item)

		y++
	}
}

func (v *ListView) drawItem(y int, fg, bg termbox.Attribute, item string) {
	x := v.StartX

	if v.EndX-v.StartX > 1 {
		termbox.SetCell(x, y, ' ', fg, bg)
		x++
	}

	for _, c := range item {
		if x >= v.EndX {
			break
		}

		termbox.SetCell(x, y, c, fg, bg)

		x++
	}

	for x < v.EndX {
		termbox.SetCell(x, y, ' ', fg, bg)
		x++
	}
}

func (v *ListView) CursorDown() {
	if v.selected >= len(v.Items)-1 {
		return
	}

	v.selected++
	for !v.isVisible(v.selected) {
		v.top++
	}
}

func (v *ListView) CursorUp() {
	if v.selected < 1 {
		return
	}

	v.selected--
	for !v.isVisible(v.selected) {
		v.top--
	}
}

func (v *ListView) isVisible(i int) bool {
	if i < v.top {
		return false
	}
	rows := v.EndY - v.StartY
	if i >= v.top+rows {
		return false
	}
	return true
}

func (v *ListView) HandleEvent(ev termbox.Event) bool {
	if ev.Type != termbox.EventKey {
		return false
	}

	switch ev.Key {
	case termbox.KeyArrowDown:
		v.CursorDown()
		return true
	case termbox.KeyArrowUp:
		v.CursorUp()
		return true
	}

	return false
}

func (v *ListView) SelectedItem() int {
	return v.selected
}

func (v *ListView) NumItems() int {
	return len(v.Items)
}
