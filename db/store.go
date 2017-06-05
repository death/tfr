package db

import (
	"database/sql"
	"log"
)

type Section struct {
	ID    int
	Label string
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

const (
	stmtAllSections = `
SELECT id, label FROM sections
`
)
