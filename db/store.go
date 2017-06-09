package db

import (
	"database/sql"
	"errors"
	"time"
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
	ModTime   time.Time
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

func (s *Store) RandomArticle() (*Article, error) {
	return s.extractArticle(stmtRandomArticle)
}

func (s *Store) LatestReadArticle() (*Article, error) {
	return s.extractArticle(stmtLatestReadArticle)
}

func (s *Store) OldestUnfinishedArticle() (*Article, error) {
	return s.extractArticle(stmtOldestUnfinishedArticle)
}

func (s *Store) extractArticle(query string) (*Article, error) {
	rows, err := s.handle.Query(query)
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
			&article.State,
			&article.ModTime)
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

const (
	stateUnread     = 0
	stateRead       = 1
	stateUnfinished = 2
)

func (s *Store) MarkAsRead(articleID int) error {
	_, err := s.handle.Exec(stmtMarkAs, stateRead, articleID)
	return err
}

func (s *Store) MarkAsUnfinished(articleID int) error {
	_, err := s.handle.Exec(stmtMarkAs, stateUnfinished, articleID)
	return err
}

const (
	stmtRandomArticle = `
SELECT id, section_id, path, size, state, mtime FROM articles
WHERE state = 0
ORDER BY random()
LIMIT 1
`
	stmtFindSection = `
SELECT label, path FROM sections
WHERE id = ?
`
	stmtMarkAs = `
UPDATE OR IGNORE articles
SET state = ?, mtime = CURRENT_TIMESTAMP
WHERE id = ?
`
	stmtLatestReadArticle = `
SELECT id, section_id, path, size, state, mtime FROM articles
WHERE state = 1
ORDER BY mtime DESC
LIMIT 1
`
	stmtOldestUnfinishedArticle = `
SELECT id, section_id, path, size, state, mtime FROM articles
WHERE state = 2
ORDER BY mtime ASC
LIMIT 1
`
)
