package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/shibukawa/configdir"
)

type Options struct {
	BaseDir string `long:"base-dir" description:"Base directory"`
	DBFile  string `long:"db-file" description:"Database file name"`
}

var (
	options Options
	parser  = flags.NewParser(&options, flags.Default)
)

func main() {
	configDirs := configdir.New("adeht", "tfr")
	configDir := configDirs.QueryFolderContainsFile("config")
	if configDir != nil {
		ini := flags.NewIniParser(parser)
		if err := ini.ParseFile(filepath.Join(configDir.Path, "config")); err != nil {
			log.Printf("config file error: %v\n", err)
		}
	}
	if options.BaseDir == "" {
		log.Fatal("base directory must be supplied via config file")
	}
	if options.DBFile == "" {
		options.DBFile = filepath.Join(options.BaseDir, "tfr.db")
	}
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
