package logseqize

import (
	"fmt"

	"github.com/gomarkdown/markdown/ast"
)

func Modify(input ast.Node) error {
	ast.WalkFunc(input, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch node.(type) {
			case *ast.Heading:
				replaceHeadingWithParagraph(node.(*ast.Heading))
			}
		}
		return ast.GoToNext
	})

	return nil
}

func replaceHeadingWithParagraph(input ast.Node) error {
	if heading, ok := input.(*ast.Heading); !ok {
		return fmt.Errorf("input is not a node")
	} else {
		para := &ast.Paragraph{Container: heading.Container}
		for i, node := range heading.GetParent().GetChildren() {
			if node == heading {
				heading.GetParent().GetChildren()[i] = para
			}
		}
	}
	return nil
}
