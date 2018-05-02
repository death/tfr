package db

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Section struct {
	ID    int
	Label string
	Path  string
}

const AnySection = -1
const AnySectionLabel = "any"

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

func (s *Store) RandomArticle(sectionID int) (*Article, error) {
	return s.extractArticle(stmtRandomArticle, sectionID)
}

func (s *Store) OldestUnfinishedArticle(sectionID int) (*Article, error) {
	return s.extractArticle(stmtOldestUnfinishedArticle, sectionID)
}

func (s *Store) extractArticle(query string, sectionID int) (*Article, error) {
	rows, err := s.handle.Query(query, sectionID)
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

func (s *Store) FindSectionByLabel(sectionLabel string) (*Section, error) {
	rows, err := s.handle.Query(stmtFindSectionByLabel, sectionLabel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		section := &Section{
			Label: sectionLabel,
		}
		err = rows.Scan(&section.ID,
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

func (s *Store) FindSectionByID(sectionID int) (*Section, error) {
	rows, err := s.handle.Query(stmtFindSectionByID, sectionID)
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

type Stats struct {
	Unread     int
	InProgress int
	Read       int
}

func (s *Store) Statistics() (map[string]*Stats, error) {
	rows, err := s.handle.Query(stmtStatistics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statsMap := make(map[string]*Stats)

	for rows.Next() {
		var section string
		var state int
		var count int
		if err := rows.Scan(&section, &state, &count); err != nil {
			return nil, err
		}
		if statsMap[section] == nil {
			statsMap[section] = &Stats{}
		}
		stats := statsMap[section]
		switch state {
		case stateUnread:
			stats.Unread = count
		case stateRead:
			stats.Read = count
		case stateUnfinished:
			stats.InProgress = count
		default:
			log.Printf("Statistics: unexpected state %d\n", state)
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return statsMap, nil
}

const (
	stmtRandomArticle = `
SELECT id, section_id, path, size, state, mtime FROM articles
WHERE state = 0 AND ? IN (-1, section_id)
ORDER BY random()
LIMIT 1
`
	stmtListSections = `
SELECT id, label, path FROM sections
`
	stmtFindSectionByLabel = `
SELECT id, path FROM sections
WHERE label = ?
`
	stmtFindSectionByID = `
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
WHERE state = 2 AND ? IN (-1, section_id)
ORDER BY mtime ASC
LIMIT 1
`
	stmtStatistics = `
SELECT label, state, COUNT(*) FROM articles
JOIN sections
WHERE sections.id = section_id
GROUP BY section_id, state
`
)
