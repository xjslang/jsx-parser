package jsxparser

import (
	"fmt"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

// ExampleParseJsxExpression demonstrates basic JSX parsing
func ExampleParseJsxExpression() {
	input := `let greeting = <h1>Hello, World!</h1>`

	l := lexer.New(input)
	p := parser.New(l)
	p.UseExpressionHandler(ParseJsxExpression)
	ast := p.ParseProgram()

	fmt.Println(ast.String())
	// Output: let greeting = React.createElement("h1", null, "Hello,World!")
}

// ExampleParseJsxExpression_selfClosing demonstrates self-closing elements
func ExampleParseJsxExpression_selfClosing() {
	input := `let image = <img src="photo.jpg" alt="A photo" />`

	l := lexer.New(input)
	p := parser.New(l)
	p.UseExpressionHandler(ParseJsxExpression)
	ast := p.ParseProgram()

	fmt.Println(ast.String())
	// Output: let image = React.createElement("img", {"src": "photo.jpg", "alt": "A photo"})
}

// ExampleParseJsxExpression_nested demonstrates nested JSX elements
func ExampleParseJsxExpression_nested() {
	input := `let card = <div><span>Nested</span></div>`

	l := lexer.New(input)
	p := parser.New(l)
	p.UseExpressionHandler(ParseJsxExpression)
	ast := p.ParseProgram()

	fmt.Println(ast.String())
	// Output: let card = React.createElement("div", null, React.createElement("span", null, "Nested"))
}

// ExampleParseJsxExpression_withAttributes demonstrates JSX with attributes
func ExampleParseJsxExpression_withAttributes() {
	input := `let container = <div className="main" id="app">Content</div>`

	l := lexer.New(input)
	p := parser.New(l)
	p.UseExpressionHandler(ParseJsxExpression)
	ast := p.ParseProgram()

	fmt.Println(ast.String())
	// Output: let container = React.createElement("div", {"className": "main", "id": "app"}, "Content")
}

// ExampleParseJsxExpression_complex demonstrates complex nested JSX with attributes
func ExampleParseJsxExpression_complex() {
	input := `let app = <div className="app"><header><h1>Title</h1></header><main><p>Content</p></main></div>`

	l := lexer.New(input)
	p := parser.New(l)
	p.UseExpressionHandler(ParseJsxExpression)
	ast := p.ParseProgram()

	fmt.Println(ast.String())
	// Output: let app = React.createElement("div", {"className": "app"}, React.createElement("header", null, React.createElement("h1", null, "Title")), React.createElement("main", null, React.createElement("p", null, "Content")))
}
