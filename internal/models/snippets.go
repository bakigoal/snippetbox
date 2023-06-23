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

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	lastInsertId := 0

	stmt := `
INSERT INTO snippetbox.snippets (title, content, created, expires)
VALUES ($1, $2, $3, $4)
RETURNING id
`
	createTime := time.Now()
	expireTime := time.Now().Add(time.Hour * time.Duration(24*expires))
	err := m.DB.QueryRow(stmt, title, content, createTime, expireTime).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	s := &Snippet{}
	stmt := `
SELECT id, title, content, created, expires FROM snippetbox.snippets 
WHERE expires > now() AND id = $1
`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) Latest10() ([]*Snippet, error) {
	stmt := `
SELECT id, title, content, created, expires FROM snippetbox.snippets 
WHERE expires > now() ORDER BY id DESC LIMIT 10
`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*Snippet
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
