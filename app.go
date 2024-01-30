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

type NotebookOutputData struct {
	ImgPng   string   `json:"image/png"`
	TextHtml []string `json:"text/html"`
}

type NotebookOutput struct {
	OutputType string `json:"output_type"`
	Text       []string
	Data       NotebookOutputData
}

type NotebookCell struct {
	CellType string `json:"cell_type"`
	Source   []string
	Outputs  []NotebookOutput
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
		// Write source
		src := strings.Join(cell.Source, "")
		switch cell.CellType {
		case "markdown":
			goldmark.Convert([]byte(src), writer)
		case "code":
			iterator, _ := lexer.Tokenise(nil, src)
			formatter.Format(writer, style, iterator)
		}

		// Write outputs
		for _, output := range cell.Outputs {
			switch output.OutputType {
			case "display_data":
				if output.Data.ImgPng != "" {
					writer.Write([]byte("\n<img alt=\"img\" src=\"data:image/png;base64,"))
					writer.Write([]byte(output.Data.ImgPng))
					writer.Write([]byte("\"/>\n"))
				}
				for _, line := range output.Data.TextHtml {
					// Gitea sanitizer doesn't like spaces in src attributes
					// This is a quick and dirty way of fixing data uri images
					line = strings.Replace(line, ";base64, ", ";base64,", 1)
					writer.Write([]byte(line))
				}
			case "stream":
				writer.Write([]byte("\n<pre>\n"))
				for _, line := range output.Text {
					writer.Write([]byte(line))
				}
				writer.Write([]byte("</pre>\n"))
			}
		}
	}

	// Finished
	writer.Flush()
}
