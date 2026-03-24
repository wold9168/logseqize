package main

import (
	"fmt"

	"github.com/wold9168/logseqize"
)

func main() {
	fmt.Print("--- markdown_foo:\n\n")
	logseqize.GetFooAst("markdown_foo.md")
	fmt.Print("--- logseq_foo:\n\n")
	logseqize.GetFooAst("logseq_foo.md")
}
