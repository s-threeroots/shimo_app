package main

import (
	"html/template"
	"io"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	funcMap := template.FuncMap{}

	t := &Template{

		templates: template.Must(template.New("t").Funcs(funcMap).ParseGlob("views/*.html")),
	}

	e := echo.New()

	e.Renderer = t

	Migrate()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/files", "files")
	e.GET("/estimation/new", CreatePage)
	e.GET("/estimation/:id/edit", EditPage)
	e.GET("/api/estimation/:id", GetEstimation)
	e.GET("/api/estimation/print/:id", Print)
	e.POST("/api/estimation/:id", SaveHandler)
	e.POST("/api/estimation/duplicate/:id", DuplicateHandler)
	e.DELETE("/api/estimation/group/:id", DeleteGroupHandler)
	e.DELETE("/api/estimation/item/:id", DeleteItemHandler)
	e.Logger.Fatal(e.Start(":" + port))
}
