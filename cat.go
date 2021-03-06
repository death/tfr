package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

type CatCommand struct {
	ArchiveFile  string `long:"archive-file" description:"The tar.gz archive file" required:"true"`
	TextfilePath string `long:"textfile-path" description:"Path of the textfile in the archive" required:"true"`
}

var cat CatCommand

func init() {
	parser.AddCommand("cat",
		"Print textfile to standard output",
		"Print textfile to standard output.",
		&cat)
}

func (c CatCommand) Execute(args []string) error {
	archive, err := os.Open(c.ArchiveFile)
	if err != nil {
		return err
	}
	defer archive.Close()

	gzReader, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		h, err := tarReader.Next()
		if err != nil {
			break
		}

		if h.Name != c.TextfilePath {
			continue
		}

		filter := &CP437Filter{Reader: tarReader}

		_, err = io.Copy(os.Stdout, filter)
		if err != nil {
			return err
		}
	}

	return nil
}
