package parser

import (
	"fmt"
	"strings"
)

type QuestionType int

const (
	FillInBlanks = iota
	DropDown
	MultipleChoice
)

func (t QuestionType) String() string {
	switch t {
	case FillInBlanks:
		return "Fill in the Blanks"
	case DropDown:
		return "Drop-Down"
	case MultipleChoice:
		return "Multiple Choice"
	default:
		return ""
	}
}

type Questions struct {
	Type    QuestionType
	Text    string
	Options []string
	Correct string
}

func (question Questions) String() string {
	optStr := ""
	if len(question.Options) > 0 {
		optStr += fmt.Sprintf("[%s]", strings.Join(question.Options, ", "))
	} else {
		optStr = "none"
	}

	return fmt.Sprintf(" - Type: %s\n   - Question: %v\n   - Options: %v\n   - Correct: %s\n\n", question.Type, question.Text, optStr, question.Correct)
}

type Worksheet struct {
	Title        string
	Instructions string
	Questions    []Questions
}

func (w Worksheet) String() string {
	var strQuestions []string
	for i, opt := range w.Questions {
		strQuestions = append(strQuestions, fmt.Sprintf("%d)%v", i+1, opt))
	}
	str := fmt.Sprintf(" - Title: %s\n - Instructions: %s\n - Questions:\n\n%v", w.Title, w.Instructions, strings.Join(strQuestions, "\n"))

	return str
}
