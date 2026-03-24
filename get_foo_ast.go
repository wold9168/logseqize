package logseqize

import (
	"fmt"
	"os"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// GetFooAst prints AST of markdown_foo.md into stdout
// Provides reference for development.
func GetFooAst(filename string) {
	md, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile Failed:", err)
		return
	}
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	fmt.Printf("-- Markdown:\n%s\n\n", md)
	ast.Print(os.Stdout, doc)
}
