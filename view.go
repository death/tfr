package main

import "github.com/nsf/termbox-go"

type View interface {
	Draw()
	HandleEvent(ev termbox.Event) bool
}
