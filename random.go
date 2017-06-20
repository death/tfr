package main

type RandomCommand struct {
}

var random RandomCommand

func init() {
	parser.AddCommand("random",
		"Read a random unread textfile",
		`Pick a random textfile for viewing; it will also be added
to the queue of unfinished textfiles.`,
		&random)
}

func (c *RandomCommand) Execute(args []string) error {
	return ViewArticle(PickRandom)
}
