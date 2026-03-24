package logseqize

import (
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

func ConvertByString(input string) (output string, err error) {
	inputBytes := []byte(input)
	output,err = Convert(inputBytes)
	return
}

func Convert(input []byte) (output string, err error) {
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	doctree := p.Parse(input)
	Modify(&doctree)
	output, err = ast.ToString(doctree), nil
	return
}
func DeConvert(input []byte) (output []byte, err error) {
	return
}
