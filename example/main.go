package main

import (
	"fmt"
	"os"

	"github.com/wold9168/logseqize"
)

func main() {
	md, err := os.ReadFile("markdown_foo.md")
	if err != nil {
		fmt.Println("ReadFile Failed:", err)
		return
	}
	fmt.Print(logseqize.Convert(md))
}
