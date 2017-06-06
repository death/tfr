package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Section struct {
	ID    int
	Label string
}

type Article struct {
	ID    int
	Label string
	Path  string
	Size  int64
	State int
}

type Store struct {
	handle *sql.DB
}

func NewStore(dbfile string) (*Store, error) {
	handle, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}
	return &Store{
		handle: handle,
	}, nil
}

func (s *Store) Close() {
	s.handle.Close()
}

func (s *Store) AllSections() []Section {
	logError := func(err error) {
		log.Printf("Store.AllSections: %v", err)
	}
	rows, err := s.handle.Query(stmtAllSections)
	if err != nil {
		logError(err)
		return nil
	}
	defer rows.Close()

	sections := make([]Section, 0)
	for rows.Next() {
		section := Section{}
		err = rows.Scan(&section.ID, &section.Label)
		if err != nil {
			logError(err)
			continue
		}
		sections = append(sections, section)
	}
	if rows.Err() != nil {
		logError(rows.Err())
	}
	return sections
}

func (s *Store) ArticlesForSection(sectionID int) []Article {
	logError := func(err error) {
		log.Printf("Store.ArticlesForSection: %v\n", err)
	}
	rows, err := s.handle.Query(stmtArticlesForSection, sectionID)
	if err != nil {
		logError(err)
		return nil
	}
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		article := Article{}
		err = rows.Scan(&article.ID, &article.Path, &article.Size, &article.State)
		if err != nil {
			logError(err)
			continue
		}
		article.Label = fmt.Sprintf("%-60s %10d", article.Path, article.Size)
		articles = append(articles, article)
	}
	if rows.Err() != nil {
		logError(rows.Err())
	}
	return articles
}

const (
	stmtAllSections = `
SELECT id, label FROM sections
`
	stmtArticlesForSection = `
SELECT id, path, size, state FROM articles
WHERE section_id = ?
`
)
