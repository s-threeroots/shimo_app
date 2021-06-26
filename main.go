package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	funcMap := template.FuncMap{}

	t := &Template{

		templates: template.Must(template.New("t").Funcs(funcMap).ParseGlob("views/*.html")),
	}

	e := echo.New()

	e.Renderer = t

	Migrate()

	e.GET("/edit/:id", func(c echo.Context) error {
		return EditHandler(c)
	})
	e.Logger.Fatal(e.Start(":8081"))
}
