package main

import (
	"github.com/bakigoal/snippetbox/internal/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			page,
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		name := filepath.Base(page)
		cache[name] = ts
	}
	return cache, nil
}
