package generator

import (
	"embed"
	"html/template"
	"io"

	"worksheetparser/internal/models"
)

//go:embed templates/*.gohtml
var templateFS embed.FS

func RenderInteractive(w io.Writer, data models.Worksheet) error {
	tmpl, err := template.ParseFS(templateFS, "templates/interactive.gohtml")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

func RenderPrintable(w io.Writer, data models.Worksheet) error {
	tmpl, err := template.ParseFS(templateFS, "templates/printable.gohtml")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
