package parser

import (
	"strings"
	"testing"
)

func TestParseExampleShape(t *testing.T) {
	ws, warnings, err := ParseFile("../../examples/example.wsp.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
	if ws.Title != "Present Simple Tense Practice" || len(ws.Questions) != 10 {
		t.Fatalf("title=%q questions=%d", ws.Title, len(ws.Questions))
	}
	q := ws.Questions[0]
	if q.Type != "multiple_choice" || q.CorrectAnswer != "does" || strings.Join(q.Options, "|") != "do|does|doing" {
		t.Errorf("first question parsed wrong: %+v", q)
	}
	if last := ws.Questions[9]; last.Type != "fill_blank" || last.Options != nil {
		t.Errorf("last question parsed wrong: %+v", last)
	}
}

func TestParseIsStateless(t *testing.T) {
	src := "@title: T\n@instructions: I\n@type: fill_blank\n@question: Q\n@correct: a\n"
	Parse(strings.NewReader(src))
	ws, _ := Parse(strings.NewReader(src))
	if len(ws.Questions) != 1 {
		t.Fatalf("second parse saw %d questions, want 1", len(ws.Questions))
	}
}

func TestParseWarnings(t *testing.T) {
	src := strings.Join([]string{
		"@title: T",
		"just some text", // not a key line
		"@nocolon",       // missing ':'
		"@colour: red",   // unknown key
		"@type: essay",   // unknown type
		"@question: Q1",  // no @correct before next question
		"@type: dropdown",
		"@question: Q2",
		"@options: a, b",
		"@correct: c",          // not in options
		"@question: Q3 at EOF", // unterminated
	}, "\n")
	ws, warnings := Parse(strings.NewReader(src))
	if len(ws.Questions) != 3 {
		t.Fatalf("got %d questions, want 3 (unterminated ones are kept)", len(ws.Questions))
	}
	if ws.Questions[2].Type != "fill_blank" {
		t.Errorf("question without @type should default to fill_blank, got %q", ws.Questions[2].Type)
	}
	for _, want := range []string{"line 2:", "line 3:", "@colour", `"essay"`, "question 1 has no @correct", `"c" is not one of the options`, "question 3 has no @correct", "missing @instructions"} {
		if !strings.Contains(strings.Join(warnings, "\n"), want) {
			t.Errorf("missing warning containing %q in:\n%s", want, strings.Join(warnings, "\n"))
		}
	}
}

func TestParseOptionsTrimmed(t *testing.T) {
	ws, _ := Parse(strings.NewReader("@type: dropdown\n@question: Q\n@options:  x ,y,, z \n@correct: y"))
	if got := strings.Join(ws.Questions[0].Options, "|"); got != "x|y|z" {
		t.Errorf("options = %q", got)
	}
}
