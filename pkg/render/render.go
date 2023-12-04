package render

import (
	"bytes"
	"fmt"
	"log"
	"myapp/pkg/config"
	"myapp/pkg/models"
	"net/http"
	"path/filepath"
	"text/template"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(response http.ResponseWriter, tmpl string, td *models.TemplateData) {

	// create template cache
	tc := app.TemplateCache

	// get requested template
	t, isOk := tc[tmpl]

	if !isOk {
		log.Fatal("Could not get template")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	_ = t.Execute(buf, td)

	// render template
	_, err := buf.WriteTo(response)

	if err != nil {
		fmt.Println("Error =: ", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCahce := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.html")

	if err != nil {
		return myCahce, err
	}

	// range through the pages
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCahce, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCahce, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCahce, err
			}
		}

		myCahce[name] = ts

	}

	return myCahce, nil

}
