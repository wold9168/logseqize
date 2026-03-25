package logseqize

import (
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md"
	"github.com/gomarkdown/markdown/parser"
)

// ConvertByString converts Standard Markdown String into Logseq-style Markdown String
func ConvertByString(input string) (output string, err error) {
	inputBytes := []byte(input)
	output, err = Convert(inputBytes)
	return
}

// Convert converts Standard Markdown []byte into Logseq-style Markdown String
func Convert(input []byte) (output string, err error) {
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	doctree := p.Parse(input)
	Modify(doctree)
	renderer := md.NewRenderer()

	output, err = string(markdown.Render(doctree, renderer)), nil
	return
}

// DeConvert haven't been implemented yet.
func DeConvert(input []byte) (output []byte, err error) {
	return
}
