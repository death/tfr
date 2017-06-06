package main

import (
	"container/list"
	"log"
	"os"

	"github.com/death/tfr/db"
	"github.com/nsf/termbox-go"
)

const (
	textfilesDir    = "/media/1984/Documents/textfiles/"
	textfilesDBFile = "/media/1984/Documents/textfiles/tfr.db"
)

var viewStack *list.List

func main() {
	if _, err := os.Stat(textfilesDBFile); os.IsNotExist(err) {
		if err := db.Build(textfilesDir, textfilesDBFile); err != nil {
			log.Fatal(err)
		}
	}

	store, err := db.NewStore(textfilesDBFile)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	viewStack = list.New()
	viewStack.PushFront(NewSectionsView(store))

	draw()
	eventLoop()
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	currentView().Draw()

	termbox.Flush()
}

func eventLoop() {
	for {
		ev := termbox.PollEvent()

		if ev.Type != termbox.EventKey {
			continue
		}

		if ev.Key == termbox.KeyEsc {
			viewStack.Remove(viewStack.Front())
			if viewStack.Front() == nil {
				return
			}
			draw()
		} else {
			handled, nextView := currentView().HandleEvent(ev)
			if handled {
				if nextView != nil {
					viewStack.PushFront(nextView)
				}
				draw()
			}
		}
	}
}

func currentView() View {
	element := viewStack.Front()
	if element == nil {
		log.Fatal("No current view")
	}
	return element.Value.(View)
}
