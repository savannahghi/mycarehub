package html

import (
	"embed"
	"html/template"
	"io"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

//go:embed *
var files embed.FS

func parse(file string) *template.Template {
	return template.Must(
		template.New("layout.html").ParseFS(files, "layout.html", file))
}

type LoginParams struct {
	Title        string
	HasError     bool
	ErrorMessage string
}

func ServeLoginPage(w io.Writer, p LoginParams) error {
	tmpl := parse("login.html")
	return tmpl.Execute(w, p)
}

type ProgramChooserParams struct {
	Title             string
	AvailablePrograms []*domain.Program
	HasError          bool
	ErrorMessage      string
}

func ServeProgramChooserPage(w io.Writer, p ProgramChooserParams) error {
	tmpl := parse("program_chooser.html")
	return tmpl.Execute(w, p)
}

type FacilityChooserParams struct {
	Title               string
	AvailableFacilities []*domain.Facility
	HasError            bool
	ErrorMessage        string
}

func ServeFacilityChooserPage(w io.Writer, p FacilityChooserParams) error {
	tmpl := parse("facility_chooser.html")
	return tmpl.Execute(w, p)
}
