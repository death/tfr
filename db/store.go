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

func (s *Store) ListSections() ([]*Section, error) {
	rows, err := s.handle.Query(stmtListSections)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sections := make([]*Section, 0)

	for rows.Next() {
		section := &Section{}
		err = rows.Scan(&section.ID,
			&section.Label,
			&section.Path)
		if err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return sections, nil
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

func (s *Store) SetLatestViewedArticleID(articleID int) error {
	_, err := s.handle.Exec(stmtSetLatestViewedArticleID, articleID)
	return err
}

const (
	NoArticle = -1
)

func (s *Store) LatestViewedArticleID() (int, error) {
	rows, err := s.handle.Query(stmtLatestViewedArticleID)
	if err != nil {
		return NoArticle, err
	}
	defer rows.Close()

	var articleID int

	if !rows.Next() {
		return NoArticle, errors.New("no latest viewed article id column?")
	}

	if err := rows.Scan(&articleID); err != nil {
		return NoArticle, err
	}

	return articleID, nil
}

const (
	stmtRandomArticle = `
SELECT id, section_id, path, size, state, mtime FROM articles
WHERE state = 0
ORDER BY random()
LIMIT 1
`
	stmtListSections = `
SELECT id, label, path FROM sections
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
	stmtLatestViewedArticleID = `
SELECT value FROM globalstate
WHERE key = 'latest_viewed_article_id'
`
	stmtSetLatestViewedArticleID = `
UPDATE globalstate
SET value = ?
WHERE key = 'latest_viewed_article_id'
`
	stmtOldestUnfinishedArticle = `
SELECT id, section_id, path, size, state, mtime FROM articles
WHERE state = 2
ORDER BY mtime ASC
LIMIT 1
`
)
