package Parser

import (
	gr "github.com/PlayerR9/LyneParser/Grammar"

	"strings"

	tr "github.com/PlayerR9/MyGoLib/CustomData/Tree"
)

/////////////////////////////////////////////////////////////

const (
	Indentation string = "|  "
)

/*
// findInvalidTokenIndex finds the index of the first invalid token in the data.
// The function returns -1 if no invalid token is found.
//
// Parameters:
//
//   - branch: The branch of tokens to search for.
//   - data: The data to search in.
//
// Returns:
//
//   - int: The index of the first invalid token.
func findInvalidTokenIndex(branch []gr.LeafToken, data []byte) int {
	pos := 0

	for _, token := range branch {
		b := []byte(token.Data)

		startIndex := slext.FindSubsliceFrom(data, b, pos)
		if startIndex == -1 {
			return -1
		}

		pos += startIndex + len(token.Data)
	}

	if pos >= len(data) {
		return -1
	}

	return pos
}
*/

// FormatSyntaxError formats a syntax error in the data.
// The function returns a string with the faulty line and a caret pointing to the invalid token.
//
// Parameters:
//
//   - branch: The branch of tokens to search for.
//   - data: The data to search in.
//
// Returns:
//
//   - string: The formatted syntax error.
//
// Example:
//
//	data := []byte("Hello, word!")
//	branch := []gr.LeafToken{
//	  {Data: "Hello", At: 0},
//	  {Data: ",", At: 6},
//	  {Data: "word", At: 8}, // Invalid token (expected "world")
//	  {Data: "!", At: 12},
//	}
//
//	fmt.Println(FormatSyntaxError(branch, data))
//
// Output:
//
//	Hello, word!
//	       ^
func FormatSyntaxError(root gr.Tokener, data []byte) string {
	return TokenerString(root)

	/*
		root.Data

		firstInvalid := findInvalidTokenIndex(branch, data)
		if firstInvalid == -1 {
			return string(data)
		}

		var builder strings.Builder

		before := data[:firstInvalid]
		after := data[firstInvalid:]

		// Write all lines before the one containing the invalid token

		beforeLines := sext.ByteSplitter(before, '\n')

		if len(beforeLines) > 1 {
			builder.WriteString(sext.JoinBytes(beforeLines[:len(beforeLines)-1], '\n'))
			builder.WriteRune('\n')
		}

		// Write the faulty line
		faultyLine := beforeLines[len(beforeLines)-1]
		afterLines := sext.ByteSplitter(after, '\n')

		builder.WriteString(string(faultyLine))

		if len(afterLines) > 0 {
			builder.WriteString(string(afterLines[0]))
		}

		builder.WriteRune('\n')

		// Write the caret
		builder.WriteString(strings.Repeat(" ", len(faultyLine)))
		builder.WriteRune('^')
		builder.WriteRune('\n')

		if len(afterLines) > 1 {
			builder.WriteString(sext.JoinBytes(afterLines[1:], '\n'))
		}

		return builder.String()
	*/
}

func TokenerString(root gr.Tokener) string {
	type helper struct {
		indent string
		root   gr.Tokener
	}

	var builder strings.Builder

	t := tr.NewTraverser(
		func(f helper) error {
			builder.WriteString(f.indent)
			builder.WriteString(f.root.GetID())

			switch root := f.root.(type) {
			case *gr.LeafToken:
				builder.WriteString(" -> ")
				builder.WriteString(root.Data)
			case *gr.NonLeafToken:
				builder.WriteString(" :")
			}

			builder.WriteString("\n")

			return nil
		},
		func(f helper) ([]helper, error) {
			switch root := f.root.(type) {
			case *gr.LeafToken:
				return nil, nil
			case *gr.NonLeafToken:
				if len(root.Data) == 0 {
					return nil, nil
				}

				children := make([]helper, 0, len(root.Data))

				newIndent := f.indent + Indentation

				for _, child := range root.Data {
					children = append(children, helper{
						indent: newIndent,
						root:   child,
					})
				}

				return children, nil
			}

			return nil, nil
		},
	)

	rootNode := helper{
		indent: "",
		root:   root,
	}

	err := t.DFS(rootNode)
	if err != nil {
		return ""
	}

	return builder.String()
}
