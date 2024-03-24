package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type templates struct {
	templates *template.Template
}

func (t *templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplates() *templates {
	return &templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type IndexModel struct {
	Count int
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newTemplates()

	indexModel := IndexModel{Count: 0}
	e.GET("/", func(c echo.Context) error {
		indexModel.Count++
		return c.Render(200, "index", indexModel)
	})

	e.Logger.Fatal(e.Start(":8081"))
}
