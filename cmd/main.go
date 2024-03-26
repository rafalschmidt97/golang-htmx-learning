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
	FormName         string
	FormAddress      string
	Contacts         []Contact
}

func (indexModel *IndexModel) addContact(name, address string) (contact Contact, error error) {
	if indexModel.contactExists(address) {
		indexModel.FormErrorMessage = "Contact with this address already exists"
		return Contact{}, errors.New("exists")
	}

	newContact := Contact{Name: name, Address: address}
	indexModel.Contacts = append(indexModel.Contacts, newContact)
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
		indexModel.FormName = c.FormValue("name")
		indexModel.FormAddress = c.FormValue("address")

		if indexModel.FormName == "" || indexModel.FormAddress == "" {
			indexModel.FormErrorMessage = "Fields must be filled"
			return c.Render(400, "form", indexModel)
		}

		co, ce := indexModel.addContact(indexModel.FormName, indexModel.FormAddress)

		if ce != nil {
			return c.Render(409, "form", indexModel)
		}

		indexModel.FormErrorMessage = ""
		indexModel.FormAddress = ""
		indexModel.FormName = ""

		c.Render(200, "form", indexModel)       // clear
		return c.Render(200, "oob-contact", co) // append
	})

	e.Logger.Fatal(e.Start(":8081"))
}
