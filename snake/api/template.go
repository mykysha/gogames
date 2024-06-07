package api

import (
	"fmt"
	"html/template"
	"io"

	"github.com/mykysha/gogames/snake/domain"
	"github.com/mykysha/gogames/snake/pkg/converters"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any) error {
	if err := t.templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func newTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("assets/views/*.html")),
	}
}

type PageData struct {
	Screen template.HTML
}

func makeInitPageData() PageData {
	screen := template.HTML("The game will start soon")

	return PageData{
		Screen: screen,
	}
}

type IndexPage struct {
	Data PageData
}

func newIndexPage() *IndexPage {
	return &IndexPage{
		Data: makeInitPageData(),
	}
}

func (p *IndexPage) UpdateScreen(screen [][]domain.Cell) {
	p.Data.Screen = converters.ConvertScreenToHTML(screen)
}
