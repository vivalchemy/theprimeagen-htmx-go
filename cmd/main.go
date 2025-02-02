package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

// type Count struct {
// 	Count int
// }

type Contact struct {
	Name  string
	Email string
}

func newContact(name string, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

type Data struct {
	Contacts []Contact
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

type Page struct {
	Data Data
	Form FormData
}

func (d *Data) hasEmail(email string) bool {
	for _, c := range d.Contacts {
		if c.Email == email {
			return true
		}
	}
	return false
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("Jhon Doe", "jhon@doe.com"),
			newContact("a", "a@b.c"),
			newContact("Clara Brown", "clara@brown.com"),
		},
	}
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

func newPageData() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// count := Count{Count: 0}
	// data := newData()
	page := newPageData()
	e.Renderer = newTemplate()

	// return c.Render(http.StatusOK, "index", count) used to test count
	e.GET("/", func(c echo.Context) error {
		// return c.Render(http.StatusOK, "index", data)
		return c.Render(http.StatusOK, "index", page)
	})

	// e.POST("/count", func(c echo.Context) error {
	// 	count.Count++
	// 	return c.Render(http.StatusOK, "count", count)
	// })

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already in use"
			return c.Render(http.StatusUnprocessableEntity, "form", formData)
		}
		contact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)
		c.Render(http.StatusOK, "form", newFormData())
		fmt.Println(contact)
		return c.Render(http.StatusOK, "oob-contact", contact)
	})
	e.Logger.Fatal(e.Start(":3000"))
}
