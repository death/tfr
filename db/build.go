package db

import (
	"archive/tar"
	"compress/gzip"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Build(dir string, dbfile string) error {
	os.Remove(dbfile)

	log.Printf("Building database file '%s'", dbfile)

	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(stmtCreateTables)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = addAllSections(dir, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func addAllSections(dir string, tx *sql.Tx) error {
	addSection, err := tx.Prepare(stmtAddSection)
	if err != nil {
		return err
	}
	defer addSection.Close()

	addArticle, err := tx.Prepare(stmtAddArticle)
	if err != nil {
		return err
	}
	defer addArticle.Close()

	sectionID := 0
	articleID := 0
	mapSections(dir, func(sectionPath string, sectionLabel string) {
		if sectionLabel == AnySectionLabel {
			log.Printf("Skipping section with reserved name '%s'", sectionLabel)
			return
		}

		_, err = addSection.Exec(sectionID, sectionLabel, sectionPath)
		if err != nil {
			log.Printf("Skipping section file '%s'", sectionPath)
			return
		}
		err = mapArticles(sectionPath, func(articlePath string, articleSize int64) {
			_, err = addArticle.Exec(articleID, sectionID, articlePath, articleSize)
			if err != nil {
				log.Printf("%s: Skipping article '%s'", sectionLabel, articlePath)
				return
			}
			articleID++
		})
		if err != nil {
			log.Printf("%s: Error mapping articles: %v", sectionLabel, err)
		}
		sectionID++
	})

	return nil
}

func mapSections(dir string, fn func(path string, label string)) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".gz" {
			return nil
		}
		if filepath.Ext(strings.TrimSuffix(path, ".gz")) != ".tar" {
			return nil
		}
		fn(path, sectionLabel(path))
		return nil
	})
}

func sectionLabel(path string) string {
	return strings.Title(strings.TrimSuffix(filepath.Base(path), ".tar.gz"))
}

func mapArticles(sectionPath string, fn func(path string, size int64)) error {
	sectionFile, err := os.Open(sectionPath)
	if err != nil {
		return err
	}
	defer sectionFile.Close()

	gzReader, err := gzip.NewReader(sectionFile)
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

		fn(h.Name, h.Size)
	}

	return nil
}

const (
	stmtCreateTables = `
CREATE TABLE sections (
    id INTEGER NOT NULL PRIMARY KEY,
    label TEXT NOT NULL,
    path TEXT NOT NULL
);

CREATE TABLE articles (
    id INTEGER NOT NULL PRIMARY KEY,
    section_id INTEGER NOT NULL,
    path TEXT NOT NULL,
    size INTEGER NOT NULL,
    state INTEGER NOT NULL DEFAULT 0,
    mtime TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE globalstate (
    key TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);

INSERT INTO globalstate (key, value) VALUES ('latest_viewed_article_id', -1);
`
	stmtAddSection = `
INSERT INTO sections (id, label, path) VALUES (?, ?, ?);
`
	stmtAddArticle = `
INSERT INTO articles (id, section_id, path, size) VALUES (?, ?, ?, ?);
`
)
