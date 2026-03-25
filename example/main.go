package main

import (
	"fmt"
	"os"

	"github.com/wold9168/logseqize"
)

func main() {
	// fmt.Print("--- markdown_foo:\n\n")
	// logseqize.GetFooAst("markdown_foo.md")
	// fmt.Print("--- logseq_foo:\n\n")
	// logseqize.GetFooAst("logseq_foo.md")
	md, err := os.ReadFile("markdown_foo.md")
	if err != nil {
		fmt.Println("ReadFile Failed:", err)
		return
	}
	fmt.Print(logseqize.Convert(md))
}
