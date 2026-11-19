//go:build js && wasm

// Exposes the parser and generator to the browser editor in web/.
// JS: wsp.render(text) -> {title, interactive, printable, warnings, questions}
package main

import (
	"strings"
	"syscall/js"

	"github.com/bugracdnc/worksheetparser/internal/generator"
	"github.com/bugracdnc/worksheetparser/internal/parser"
)

func render(_ js.Value, args []js.Value) any {
	ws, warnings := parser.Parse(strings.NewReader(args[0].String()))
	var interactive, printable strings.Builder
	if err := generator.RenderInteractive(&interactive, ws); err != nil {
		warnings = append(warnings, err.Error())
	}
	if err := generator.RenderPrintable(&printable, ws); err != nil {
		warnings = append(warnings, err.Error())
	}
	w := make([]any, len(warnings))
	for i, s := range warnings {
		w[i] = s
	}
	return map[string]any{
		"title":       ws.Title,
		"interactive": interactive.String(),
		"printable":   printable.String(),
		"warnings":    w,
		"questions":   len(ws.Questions),
	}
}

func main() {
	js.Global().Set("wsp", map[string]any{"render": js.FuncOf(render)})
	select {}
}
