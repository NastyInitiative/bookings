package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/NastyInitiative/bookings/internal/config"
	"github.com/NastyInitiative/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}
var app *config.AppConfig

func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(tData *models.TemplateData, r *http.Request) *models.TemplateData {
	tData.CSRFToken = nosurf.Token(r)
	return tData
}

//RenderTemplate - renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, tData *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	tData = AddDefaultData(tData, r)
	_ = t.Execute(buf, tData)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing to browser")
	}

}

/*Creates a template cache as a map*/
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
