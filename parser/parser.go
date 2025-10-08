package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"worksheetparser/debug"
)

var worksheet = Worksheet{}
var question = Questions{}

func indexOf(str string, char string) int {
	for i, v := range str {
		if string(v) == char {
			return i
		}
	}

	return -1
}

func parseSetting(setting string, value string) {
	switch setting {
	case "title":
		worksheet.Title = value
	case "instructions":
		worksheet.Instructions = value
	case "type":
		switch value {
		case "multiple_choice":
			question.Type = MultipleChoice
		case "fill_blank":
			question.Type = FillInBlanks
		default:
			question.Type = DropDown
		}
	case "question":
		question.Text = value
	case "options":
		options := strings.Split(value, ",")
		for i, opt := range options {
			options[i] = strings.Trim(opt, " ")
		}
		question.Options = options
	case "correct":
		question.Correct = strings.Trim(value, " ")
		worksheet.Questions = append(worksheet.Questions, question)
		question = Questions{}
	}

}

func parseLine(lineNo int, line string) {
	if len(line) < 1 {
		debug.LogPrintf("Empty line on %d -- skipped\n\n", lineNo)
		return
	}
	switch line[0] {
	case '#':
		debug.LogPrintf("Comment on line %d -- skipped\n\n", lineNo)
	case '@':
		debug.LogPrintf("Found line with markdown... : '%s'\n", line)
		setting := line[1:indexOf(line, ":")]
		value := strings.Trim(line[indexOf(line, ":")+1:], " ")
		debug.LogPrintf("setting: '%s'\nvalue: '%s'\n\n", setting, value)
		parseSetting(setting, value)
	default:
		debug.LogPrintf("Invalid line on %d -- '%s' -- skipped\n\n", lineNo, line)
		return
	}
}

func Parse(filename string) (Worksheet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Worksheet{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		parseLine(lineNo, line)

	}

	if worksheet.Title == "" {
		fmt.Println(" ! Warning: parameter missing -- title")
	}
	if worksheet.Instructions == "" {
		fmt.Println(" ! Warning: parameter missing -- instructions")
	}
	if len(worksheet.Questions) < 1 {
		fmt.Println(" ! Warning: parameter missing -- questions")
	}

	return worksheet, nil
}
