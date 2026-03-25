package logseqize

import (
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

func Modify(input ast.Node) error {
	// First pass: transform individual nodes
	ast.WalkFunc(input, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n := node.(type) {
			case *ast.BlockQuote:
				transformBlockQuote(n)
			case *ast.List:
				transformList(n)
			}
		}
		return ast.GoToNext
	})

	// Transform headings (need to replace node type)
	transformHeadings(input)

	// Transform blockquotes to paragraphs
	transformBlockQuotesToParagraphs(input)

	// Transform horizontal rules
	transformHorizontalRules(input)

	// Transform tables to paragraphs
	transformTables(input)

	// Second pass: wrap top-level nodes in a List
	wrapTopLevelInList(input)

	return nil
}

func transformHeadings(doc ast.Node) {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if heading, ok := node.(*ast.Heading); ok {
				// Prepend '#' characters based on level to the first text child
				for _, child := range heading.GetChildren() {
					if text, ok := child.(*ast.Text); ok {
						text.Literal = []byte(strings.Repeat("#", heading.Level) + " " + string(text.Literal))
						break
					}
				}
				// Convert to Paragraph
				para := &ast.Paragraph{}
				para.Children = heading.Children
				for _, child := range para.Children {
					child.SetParent(para)
				}
				// Replace heading with paragraph in parent
				parent := heading.GetParent()
				if parent != nil {
					children := parent.GetChildren()
					for i, child := range children {
						if child == heading {
							children[i] = para
							para.SetParent(parent)
							break
						}
					}
				}
				return ast.GoToNext
			}
		}
		return ast.GoToNext
	})
}

func transformBlockQuote(blockquote *ast.BlockQuote) {
	// Prepend '>' to text in blockquote's first paragraph
	for _, child := range blockquote.GetChildren() {
		if para, ok := child.(*ast.Paragraph); ok {
			for _, textNode := range para.GetChildren() {
				if text, ok := textNode.(*ast.Text); ok {
					text.Literal = []byte("> " + string(text.Literal))
					return
				}
			}
		}
	}
}

func transformBlockQuotesToParagraphs(doc ast.Node) {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if blockquote, ok := node.(*ast.BlockQuote); ok {
				// Get the first paragraph's text
				var textContent []byte
				for _, child := range blockquote.GetChildren() {
					if para, ok := child.(*ast.Paragraph); ok {
						for _, textNode := range para.GetChildren() {
							if text, ok := textNode.(*ast.Text); ok {
								textContent = text.Literal
								break
							}
						}
						break
					}
				}
				// Create new paragraph with the text
				para := &ast.Paragraph{}
				text := &ast.Text{}
				text.Literal = textContent
				para.Children = []ast.Node{text}
				text.SetParent(para)
				// Replace blockquote with paragraph
				parent := blockquote.GetParent()
				if parent != nil {
					children := parent.GetChildren()
					for i, child := range children {
						if child == blockquote {
							children[i] = para
							para.SetParent(parent)
							break
						}
					}
				}
				return ast.GoToNext
			}
		}
		return ast.GoToNext
	})
}

func transformHorizontalRules(doc ast.Node) {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if hr, ok := node.(*ast.HorizontalRule); ok {
				// Replace HorizontalRule with Paragraph containing "---"
				para := &ast.Paragraph{}
				text := &ast.Text{}
				text.Literal = []byte("---")
				para.Children = []ast.Node{text}
				text.SetParent(para)
				// Replace hr with paragraph
				parent := hr.GetParent()
				if parent != nil {
					children := parent.GetChildren()
					for i, child := range children {
						if child == hr {
							children[i] = para
							para.SetParent(parent)
							break
						}
					}
				}
				return ast.GoToNext
			}
		}
		return ast.GoToNext
	})
}

func transformTables(doc ast.Node) {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if table, ok := node.(*ast.Table); ok {
				// Convert table to paragraph with table text representation
				textContent := renderTableToText(table)
				para := &ast.Paragraph{}
				text := &ast.Text{}
				text.Literal = textContent
				para.Children = []ast.Node{text}
				text.SetParent(para)
				// Replace table with paragraph
				parent := table.GetParent()
				if parent != nil {
					children := parent.GetChildren()
					for i, child := range children {
						if child == table {
							children[i] = para
							para.SetParent(parent)
							break
						}
					}
				}
				return ast.GoToNext
			}
		}
		return ast.GoToNext
	})
}

func renderTableToText(table *ast.Table) []byte {
	var result strings.Builder
	var rows [][]string

	// Collect all rows
	for _, child := range table.GetChildren() {
		switch rowContainer := child.(type) {
		case *ast.TableHeader:
			for _, rowNode := range rowContainer.GetChildren() {
				if row, ok := rowNode.(*ast.TableRow); ok {
					rows = append(rows, extractRowCells(row))
				}
			}
		case *ast.TableBody:
			for _, rowNode := range rowContainer.GetChildren() {
				if row, ok := rowNode.(*ast.TableRow); ok {
					rows = append(rows, extractRowCells(row))
				}
			}
		}
	}

	// Calculate column widths
	if len(rows) == 0 {
		return []byte{}
	}
	colWidths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Build table string
	for i, row := range rows {
		for j, cell := range row {
			if j < len(colWidths) {
				result.WriteString("| ")
				result.WriteString(cell)
				// Pad with spaces
				for k := len(cell); k < colWidths[j]; k++ {
					result.WriteString(" ")
				}
				result.WriteString(" ")
			}
		}
		result.WriteString("|\n")
		// Add separator after header
		if i == 0 {
			for j := range colWidths {
				result.WriteString("| ")
				for k := 0; k < colWidths[j]; k++ {
					result.WriteString("-")
				}
				result.WriteString(" ")
			}
			result.WriteString("|\n")
		}
	}

	return []byte(result.String())
}

func extractRowCells(row *ast.TableRow) []string {
	var cells []string
	for _, cellNode := range row.GetChildren() {
		if cell, ok := cellNode.(*ast.TableCell); ok {
			var text string
			for _, textNode := range cell.GetChildren() {
				if t, ok := textNode.(*ast.Text); ok {
					text += string(t.Literal)
				}
			}
			cells = append(cells, text)
		}
	}
	return cells
}

func transformList(list *ast.List) {
	// Transform task list items: [x] -> DONE, [ ] -> LATER
	for _, item := range list.GetChildren() {
		if listItem, ok := item.(*ast.ListItem); ok {
			transformTaskListItem(listItem)
		}
	}
}

func transformTaskListItem(listItem *ast.ListItem) {
	// Look for task pattern in first paragraph
	for _, child := range listItem.GetChildren() {
		if para, ok := child.(*ast.Paragraph); ok {
			for _, textNode := range para.GetChildren() {
				if text, ok := textNode.(*ast.Text); ok {
					literal := string(text.Literal)
					// Check for task markers
					if strings.HasPrefix(literal, "[x] ") || strings.HasPrefix(literal, "[X] ") {
						text.Literal = []byte("DONE " + strings.TrimPrefix(strings.TrimPrefix(literal, "[x] "), "[X] "))
						return
					} else if strings.HasPrefix(literal, "[ ] ") {
						text.Literal = []byte("LATER " + strings.TrimPrefix(literal, "[ ] "))
						return
					}
				}
			}
		}
	}
}

// isTaskList checks if a list contains task items (checkboxes or converted DONE/LATER)
func isTaskList(list *ast.List) bool {
	for _, item := range list.GetChildren() {
		if listItem, ok := item.(*ast.ListItem); ok {
			for _, child := range listItem.GetChildren() {
				if para, ok := child.(*ast.Paragraph); ok {
					for _, textNode := range para.GetChildren() {
						if text, ok := textNode.(*ast.Text); ok {
							literal := string(text.Literal)
							// Check for both original task markers and converted DONE/LATER
							if strings.HasPrefix(literal, "[x] ") ||
								strings.HasPrefix(literal, "[X] ") ||
								strings.HasPrefix(literal, "[ ] ") ||
								strings.HasPrefix(literal, "DONE ") ||
								strings.HasPrefix(literal, "LATER ") {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

func wrapTopLevelInList(doc ast.Node) {
	// Get document children (top-level nodes)
	var topNodes []ast.Node
	for _, child := range doc.GetChildren() {
		topNodes = append(topNodes, child)
	}

	if len(topNodes) == 0 {
		return
	}

	// Create a new tight unordered list
	list := &ast.List{
		Container: ast.Container{},
		Tight:     true,
	}

	// Build list items, handling nested lists specially
	var listItems []ast.Node
	for i, node := range topNodes {
		// Check if this is a task list - flatten it into individual items
		if nestedList, ok := node.(*ast.List); ok && isTaskList(nestedList) {
			// Flatten task list: each item becomes a top-level ListItem
			for j, item := range nestedList.GetChildren() {
				if taskItem, ok := item.(*ast.ListItem); ok {
					// Set appropriate flags
					if i == 0 && j == 0 {
						taskItem.ListFlags = ast.ListItemBeginningOfList
					}
					listItems = append(listItems, taskItem)
				}
			}
			continue
		}

		// Check if this is a list that should be nested under previous item
		// Only nest if previous node is a Paragraph with "List:" text
		if nestedList, ok := node.(*ast.List); ok && i > 0 && len(listItems) > 0 {
			prevItem := listItems[len(listItems)-1]
			if prevListItem, ok := prevItem.(*ast.ListItem); ok {
				// Check if previous item contains a Paragraph with "List:" text
				if shouldNestListUnder(prevListItem) {
					// Append the nested list to the previous item
					prevListItem.Children = append(prevListItem.Children, nestedList)
					nestedList.SetParent(prevListItem)
					continue
				}
			}
		}

		listItem := &ast.ListItem{
			Container: ast.Container{},
		}

		// Set flags for first item
		if i == 0 && len(listItems) == 0 {
			listItem.ListFlags = ast.ListItemBeginningOfList
		}

		// Special handling for CodeBlock: wrap in Paragraph first
		if codeblock, ok := node.(*ast.CodeBlock); ok {
			para := &ast.Paragraph{}
			beforeText := &ast.Text{Leaf: ast.Leaf{Literal: []byte{}}}
			afterText := &ast.Text{Leaf: ast.Leaf{Literal: []byte{}}}
			para.Children = []ast.Node{beforeText, codeblock, afterText}
			beforeText.SetParent(para)
			codeblock.SetParent(para)
			afterText.SetParent(para)
			listItem.Children = []ast.Node{para}
			para.SetParent(listItem)
		} else {
			// Append the node as child
			listItem.Children = []ast.Node{node}
			node.SetParent(listItem)
		}
		listItems = append(listItems, listItem)
	}

	// Update flags for actual last item
	if len(listItems) > 0 {
		if lastItem, ok := listItems[len(listItems)-1].(*ast.ListItem); ok {
			lastItem.ListFlags |= ast.ListItemEndOfList
		}
	}

	list.Children = listItems
	for _, item := range listItems {
		item.SetParent(list)
	}

	// Replace document children with the new list
	if document, ok := doc.(*ast.Document); ok {
		document.Children = []ast.Node{list}
		list.SetParent(document)
	}
}

// shouldNestListUnder checks if a ListItem should have a list nested under it
// Returns true if the ListItem contains a Paragraph with text ending in "List:"
func shouldNestListUnder(listItem *ast.ListItem) bool {
	for _, child := range listItem.GetChildren() {
		if para, ok := child.(*ast.Paragraph); ok {
			for _, textNode := range para.GetChildren() {
				if text, ok := textNode.(*ast.Text); ok {
					literal := string(text.Literal)
					if strings.HasSuffix(literal, "List:") {
						return true
					}
				}
			}
		}
	}
	return false
}
