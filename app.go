package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"

	chroma_html "github.com/alecthomas/chroma/formatters/html"
	chroma_lexers "github.com/alecthomas/chroma/lexers"
	chroma_styles "github.com/alecthomas/chroma/styles"
	"github.com/yuin/goldmark"
)

type NotebookCell struct {
	Cell_Type string
	Source    []string
}

type Notebook struct {
	Cells []NotebookCell
}

func main() {
	// Chroma source formatting
	lexer := chroma_lexers.Get("python")
	style := chroma_styles.Get("monokai")
	formatter := chroma_html.New(chroma_html.WithClasses(true), chroma_html.WithLineNumbers(true))

	// Read from stdin
	reader := bufio.NewReader(os.Stdin)
	input, _ := io.ReadAll(reader)

	// Convert JSON bytes to object
	var notebook1 Notebook
	json.Unmarshal(input, &notebook1)

	// Write output to stdout
	writer := bufio.NewWriter(os.Stdout)

	// Loop through cells
	for _, cell := range notebook1.Cells {
		src := strings.Join(cell.Source, "")
		switch cell.Cell_Type {
		case "markdown":
			goldmark.Convert([]byte(src), writer)
		case "code":
			iterator, _ := lexer.Tokenise(nil, src)
			formatter.Format(writer, style, iterator)
		}
	}

	// Finished
	writer.Flush()
}
