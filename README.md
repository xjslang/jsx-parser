# JSX Parser Plugin for XJS

This is a plugin for the XJS transpiler that adds support for JSX syntax, transforming it into standard JavaScript `React.createElement` calls.

## Implemented Features

### ✅ Basic JSX Elements
```jsx
<div>Hello, World!</div>
// → React.createElement("div", null, "Hello,World!")
```

### ✅ Self-Closing Elements
```jsx
<img />
<br />
// → React.createElement("img", null)
// → React.createElement("br", null)
```

### ✅ Nested Elements
```jsx
<div>
  <span>Nested content</span>
</div>
// → React.createElement("div", null, React.createElement("span", null, "Nestedcontent"))
```

### ✅ JSX Attributes
```jsx
<div className="container" id="main">Content</div>
// → React.createElement("div", {"className": "container", "id": "main"}, "Content")
```

### ✅ Complex Elements
```jsx
<div className="main">
  <h1>Title</h1>
  <p>Paragraph</p>
</div>
// → React.createElement("div", {"className": "main"}, React.createElement("h1", null, "Title"), React.createElement("p", null, "Paragraph"))
```

## Usage

```go
package main

import (
    "github.com/xjslang/xjs/lexer"
    "github.com/xjslang/xjs/parser"
    jsxparser "github.com/xjslang/jsx-parser"
)

func main() {
    input := `let component = <div className="app">Hello, JSX!</div>`
    
    // Create lexer and parser
    l := lexer.New(input)
    p := parser.New(l)
    
    // Register the JSX middleware
    p.UseExpressionHandler(jsxparser.ParseJsxExpression)
    
    // Parse the program
    ast := p.ParseProgram()
    
    // Get the transpiled JavaScript code
    output := ast.String()
    // → "let component = React.createElement("div", {"className": "app"}, "Hello,JSX!")"
}
```

## Examples

For detailed usage examples, check the **Example functions** included in the code:

```bash
# See all available examples
go test -run Example -v

# Run a specific example
go test -run ExampleParseJsxExpression_nested -v
```

The examples cover:
- **`ExampleParseJsxExpression`**: Basic usage
- **`ExampleParseJsxExpression_selfClosing`**: Self-closing elements  
- **`ExampleParseJsxExpression_nested`**: Nested elements
- **`ExampleParseJsxExpression_withAttributes`**: JSX with attributes
- **`ExampleParseJsxExpression_complex`**: Complex cases

## Architecture

The plugin uses the XJS middleware pattern:

1. **Middleware Pattern**: Registers as an `ExpressionHandler` that intercepts expression parsing
2. **Token Recognition**: Detects the `<` + `IDENT` sequence to start JSX parsing
3. **Recursive Parsing**: Handles nested elements recursively
4. **Fallback**: If not valid JSX, passes control to the next parser in the chain

### Main Components

- **`JsxExpression`**: Represents a JSX element in the AST
- **`JsxAttribute`**: Represents a JSX element attribute  
- **`JsxText`**: Represents text content inside JSX elements
- **`ParseJsxExpression`**: Middleware function that handles parsing

## Current Limitations

- **JavaScript Expressions**: Does not support `{expression}` inside JSX
- **Fragments**: Does not support `<>...</>` (React Fragments)
- **Components**: Only supports native HTML elements (lowercase tags)
- **Complex Attributes**: Only supports attributes with string literal values

## Upcoming Features

- [ ] Support for JavaScript expressions `{variable}`
- [ ] React Fragments `<>...</>`
- [ ] React Components (PascalCase)
- [ ] Attributes with complex expressions
- [ ] HTML tag name validation
- [ ] Better whitespace handling in text

## Testing

```bash
# Run all tests
go test -v

# Run only unit tests
go test -run TestJsxParser -v

# Run only examples
go test -run Example -v

# View documentation with examples
go doc .
```

The plugin includes:
- **Comprehensive unit tests** covering all use cases
- **Example functions** serving as executable documentation and usage examples

