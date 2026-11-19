// Package parser reads the .wsp.txt worksheet format:
//
//	@title: ...            worksheet title
//	@instructions: ...     shown under the title
//	@type: multiple_choice | dropdown | fill_blank
//	@question: ...
//	@options: a, b, c      (multiple_choice and dropdown)
//	@correct: b            closes the question
//
// Lines starting with # are comments; blank lines are ignored.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/bugracdnc/worksheetparser/internal/models"
)

var validTypes = []string{"multiple_choice", "dropdown", "fill_blank"}

// Parse reads a worksheet. It never fails on bad input: problems are
// returned as human-readable warnings so a live editor can show them.
func Parse(r io.Reader) (models.Worksheet, []string) {
	var (
		ws       models.Worksheet
		q        models.Questions
		warnings []string
	)
	warn := func(line int, format string, a ...any) {
		warnings = append(warnings, fmt.Sprintf("line %d: ", line)+fmt.Sprintf(format, a...))
	}
	pending := func() bool { return q.Type != "" || q.Text != "" || len(q.Options) > 0 }
	finish := func(line int) {
		if q.Type == "" {
			q.Type = "fill_blank"
		}
		if q.CorrectAnswer == "" {
			warn(line, "question %d has no @correct", len(ws.Questions)+1)
		} else if len(q.Options) > 0 && !slices.Contains(q.Options, q.CorrectAnswer) {
			warn(line, "question %d: correct answer %q is not one of the options", len(ws.Questions)+1, q.CorrectAnswer)
		}
		ws.Questions = append(ws.Questions, q)
		q = models.Questions{}
	}

	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !strings.HasPrefix(line, "@") || !ok {
			warn(lineNo, "expected '@key: value', got %q", line)
			continue
		}
		key = strings.ToLower(strings.TrimSpace(key[1:]))
		value = strings.TrimSpace(value)

		switch key {
		case "title":
			ws.Title = value
		case "instructions":
			ws.Instructions = value
		case "type", "question":
			// A new question started before the previous one got its @correct.
			if (key == "type" && pending()) || (key == "question" && q.Text != "") {
				finish(lineNo)
			}
			if key == "question" {
				q.Text = value
				break
			}
			if !slices.Contains(validTypes, value) {
				warn(lineNo, "unknown type %q (use %s)", value, strings.Join(validTypes, ", "))
			}
			q.Type = value
		case "options":
			q.Options = nil
			for _, opt := range strings.Split(value, ",") {
				if opt = strings.TrimSpace(opt); opt != "" {
					q.Options = append(q.Options, opt)
				}
			}
		case "correct":
			q.CorrectAnswer = value
			finish(lineNo)
		default:
			warn(lineNo, "unknown key @%s", key)
		}
	}
	if pending() {
		finish(lineNo)
	}

	if ws.Title == "" {
		warnings = append(warnings, "missing @title")
	}
	if ws.Instructions == "" {
		warnings = append(warnings, "missing @instructions")
	}
	if len(ws.Questions) == 0 {
		warnings = append(warnings, "no questions found")
	}
	return ws, warnings
}

// ParseFile is Parse for a file on disk.
func ParseFile(path string) (models.Worksheet, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return models.Worksheet{}, nil, err
	}
	defer f.Close()
	ws, warnings := Parse(f)
	return ws, warnings, nil
}
