package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) error {
	stmt := `
INSERT INTO snippetbox.snippets (title, content, created, expires)
VALUES ($1, $2, $3, $4)
`
	createTime := time.Now()
	expireTime := time.Now().Add(time.Hour * time.Duration(24*expires))
	_, err := m.DB.Exec(stmt, title, content, createTime, expireTime)
	if err != nil {
		return err
	}
	return nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `
SELECT id, title, content, created, expires FROM snippetbox.snippets
WHERE expires > now() AND id = ?
`
	row := m.DB.QueryRow(stmt, id)
	s := &Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) Latest10() ([]*Snippet, error) {
	return nil, nil
}
