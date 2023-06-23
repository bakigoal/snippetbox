package main

import (
	"github.com/bakigoal/snippetbox/internal/models"
	"html/template"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		ts, err := createTemplate(page)
		if err != nil {
			return nil, err
		}
		cache[filepath.Base(page)] = ts
	}
	return cache, nil
}

func createTemplate(page string) (*template.Template, error) {
	// Base template
	ts, err := template.New(page).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
	if err != nil {
		return nil, err
	}
	// Partials (nav, header, footer...)
	ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	// current Page
	ts, err = ts.ParseFiles(page)
	if err != nil {
		return nil, err
	}
	return ts, err
}

var functions = template.FuncMap{
	"humanDate": func(t time.Time) string {
		return t.Format("02-01-2006 15:04")
	},
}
