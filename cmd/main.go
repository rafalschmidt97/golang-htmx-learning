package main

import (
	"errors"
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
	FormErrorMessage string
	Contacts         []Contact
}

func (indexModel *IndexModel) addContact(name, address string) (contact Contact, error error) {
	if indexModel.contactExists(address) {
		indexModel.FormErrorMessage = "Contact with this address already exists"
		return Contact{}, errors.New("exists")
	}

	newContact := Contact{Name: name, Address: address}
	indexModel.Contacts = append(indexModel.Contacts, newContact)
	indexModel.FormErrorMessage = ""
	return newContact, nil
}

func (indexModel *IndexModel) contactExists(address string) bool {
	for _, v := range indexModel.Contacts {
		if v.Address == address {
			return true
		}
	}

	return false
}

type Contact struct {
	Name    string
	Address string
}

func main() {
	e := echo.New()
	e.Static("/", "css")
	e.Use(middleware.Logger())
	e.Renderer = newTemplates()

	indexModel := IndexModel{
		Contacts: []Contact{
			{Name: "Hello", Address: "h@w.pl"},
			{Name: "World", Address: "w@w.pl"},
		},
	}

	e.GET("/", func(c echo.Context) error {
		indexModel.FormErrorMessage = ""
		return c.Render(200, "index", indexModel)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		address := c.FormValue("address")
		_, ce := indexModel.addContact(name, address)

		if ce != nil {
			return c.Render(409, "form-error", indexModel)
		}

		return c.Render(200, "contacts", indexModel)
		// return c.Render(200, "oob-contact", co)
	})

	e.Logger.Fatal(e.Start(":8081"))
}
