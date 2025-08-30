package jsxparser

import (
	"testing"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestJsxParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic JSX element",
			input:    `let x = <div>Hello, World!</div>`,
			expected: `let x = React.createElement("div", null, "Hello,World!")`,
		},
		{
			name:     "Self-closing element",
			input:    `let y = <img />`,
			expected: `let y = React.createElement("img", null)`,
		},
		{
			name:     "Nested elements",
			input:    `let z = <div><span>Nested</span></div>`,
			expected: `let z = React.createElement("div", null, React.createElement("span", null, "Nested"))`,
		},
		{
			name:     "Element with attributes",
			input:    `let w = <div className="container">Content</div>`,
			expected: `let w = React.createElement("div", {"className": "container"}, "Content")`,
		},
		{
			name:     "Complex nested with attributes",
			input:    `let complex = <div className="main"><h1>Title</h1><p>Paragraph</p></div>`,
			expected: `let complex = React.createElement("div", {"className": "main"}, React.createElement("h1", null, "Title"), React.createElement("p", null, "Paragraph"))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			p.UseExpressionHandler(ParseJsxExpression)
			ast := p.ParseProgram()
			result := ast.String()

			if result != tt.expected {
				t.Errorf("Expected: %s\nGot: %s", tt.expected, result)
			}
		})
	}
}
