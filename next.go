package main

type NextCommand struct {
	Section string `short:"s" long:"section" default:"any"`
}

var next NextCommand

func init() {
	parser.AddCommand("next",
		"View the oldest unfinished article",
		"View the oldest unfinished article.",
		&next)
}

func (c *NextCommand) Execute(args []string) error {
	return ViewArticle(PickUnfinished, c.Section)
}
