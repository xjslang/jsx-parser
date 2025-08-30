package main

import (
	"fmt"
	"strings"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	jsxparser "github.com/xjslang/jsx-parser"
)

func main() {
	examples := []string{
		// Ejemplo básico
		`let greeting = <h1>Hello, World!</h1>`,
		
		// Elemento self-closing
		`let image = <img src="photo.jpg" alt="A photo" />`,
		
		// Elementos anidados
		`let card = <div className="card">
			<h2>Card Title</h2>
			<p>Card description here.</p>
		</div>`,
		
		// Ejemplo complejo
		`let app = <div className="app">
			<header className="header">
				<h1>My App</h1>
			</header>
			<main className="main">
				<p>Welcome to the app!</p>
			</main>
		</div>`,
	}

	fmt.Println("=== JSX to JavaScript Transpilation Examples ===\n")

	for i, input := range examples {
		fmt.Printf("Example %d:\n", i+1)
		fmt.Printf("Input JSX:\n%s\n", input)
		
		// Crear lexer y parser
		l := lexer.New(input)
		p := parser.New(l)
		
		// Registrar el middleware JSX
		p.UseExpressionHandler(jsxparser.ParseJsxExpression)
		
		// Parsear el programa
		ast := p.ParseProgram()
		
		// Obtener el código JavaScript transpilado
		output := ast.String()
		fmt.Printf("Output JavaScript:\n%s\n", output)
		fmt.Println(strings.Repeat("-", 60))
	}
}
