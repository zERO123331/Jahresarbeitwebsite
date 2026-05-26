package main

import (
	"Jahresarbeitwebsite/internal/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Updates        []Update
	LatestUpdate   Update
	SpecificUpdate Update
	User           *models.User
	ShopEntries    []*models.ShopEntry
	SpecificEntry  *models.ShopEntry
	Error          string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"./ui/html/base.gohtml",
			"./ui/html/partials/nav.gohtml",
			page,
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
