package generator

import (
	"strings"
	"testing"

	"github.com/bugracdnc/worksheetparser/internal/models"
)

func TestRenderEscapesAndIncludesQuestions(t *testing.T) {
	ws := models.Worksheet{
		Title: "Quiz <script>",
		Questions: []models.Questions{
			{Type: "multiple_choice", Text: "Pick one", Options: []string{"a", "b"}, CorrectAnswer: "b"},
			{Type: "dropdown", Text: "Choose", Options: []string{"x"}, CorrectAnswer: "x"},
			{Type: "fill_blank", Text: "Type it", CorrectAnswer: "word"},
		},
	}
	for name, render := range map[string]func(*strings.Builder) error{
		"interactive": func(b *strings.Builder) error { return RenderInteractive(b, ws) },
		"printable":   func(b *strings.Builder) error { return RenderPrintable(b, ws) },
	} {
		var b strings.Builder
		if err := render(&b); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		out := b.String()
		if strings.Contains(out, "Quiz <script>") {
			t.Errorf("%s: title not escaped", name)
		}
		for _, want := range []string{"Pick one", "Choose", "Type it"} {
			if !strings.Contains(out, want) {
				t.Errorf("%s: missing %q", name, want)
			}
		}
	}
}
