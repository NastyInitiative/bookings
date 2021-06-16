package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/NastyInitiative/bookings/internal/config"
	"github.com/NastyInitiative/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}
var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer - sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData - adds data for all templates
func AddDefaultData(tData *models.TemplateData, r *http.Request) *models.TemplateData {
	tData.CSRFToken = nosurf.Token(r)
	tData.Flash = app.Session.PopString(r.Context(), "flash")
	tData.Warning = app.Session.PopString(r.Context(), "warning")
	tData.Error = app.Session.PopString(r.Context(), "error")
	return tData
}

//Template - renders templates using html/template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, tData *models.TemplateData) error {

	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		// log.Fatal("Could not get template from template cache")
		return errors.New("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	tData = AddDefaultData(tData, r)
	_ = t.Execute(buf, tData)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing to browser")
		return err
	}
	return nil
}

/*CreateTemplateCache - Creates a template cache as a map*/
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// pages, err := filepath.Glob("./templates/*.page.tmpl")
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
