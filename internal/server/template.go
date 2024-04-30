package server

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Template struct {
	Template *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Template.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
	return &Template{
		Template: template.Must(template.ParseGlob("web/**/*.html")),
	}
}
