package models

import (
	"database/sql"
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
	return nil, nil
}

func (m *SnippetModel) Latest10() ([]*Snippet, error) {
	return nil, nil
}
