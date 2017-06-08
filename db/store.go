package db

import (
	"database/sql"
	"errors"
	"log"
)

type Section struct {
	ID    int
	Label string
	Path  string
}

type Article struct {
	ID        int
	SectionID int
	Path      string
	Size      int64
	State     int
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
		err = rows.Scan(&section.ID, &section.Label, &section.Path)
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
		article := Article{
			SectionID: sectionID,
		}
		err = rows.Scan(&article.ID, &article.Path, &article.Size, &article.State)
		if err != nil {
			logError(err)
			continue
		}
		articles = append(articles, article)
	}
	if rows.Err() != nil {
		logError(rows.Err())
	}
	return articles
}

func (s *Store) RandomArticle() (*Article, error) {
	rows, err := s.handle.Query(stmtRandomArticle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		article := &Article{}
		err = rows.Scan(&article.ID,
			&article.SectionID,
			&article.Path,
			&article.Size,
			&article.State)
		if err != nil {
			return nil, err
		}
		return article, nil
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return nil, errors.New("no article found")
}

func (s *Store) FindSection(sectionID int) (*Section, error) {
	rows, err := s.handle.Query(stmtFindSection, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		section := &Section{
			ID: sectionID,
		}
		err = rows.Scan(&section.Label,
			&section.Path)
		if err != nil {
			return nil, err
		}
		return section, nil
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return nil, errors.New("no section found")
}

func (s *Store) MarkAsRead(articleID int) error {
	_, err := s.handle.Exec(stmtMarkAsRead, articleID)
	return err
}

const (
	stmtAllSections = `
SELECT id, label, path FROM sections
`
	stmtArticlesForSection = `
SELECT id, path, size, state FROM articles
WHERE section_id = ?
`
	stmtRandomArticle = `
SELECT id, section_id, path, size, state FROM articles
WHERE state = 0
ORDER BY random()
LIMIT 1
`
	stmtFindSection = `
SELECT label, path FROM sections
WHERE id = ?
`
	stmtMarkAsRead = `
UPDATE OR IGNORE articles
SET state = 1
WHERE id = ?
`
)
