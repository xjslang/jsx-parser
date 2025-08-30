# JSX Parser Plugin para XJS

Este es un plugin para el transpilador XJS que añade soporte para sintaxis JSX, transformándola a llamadas `React.createElement` estándar de JavaScript.

## Características Implementadas

### ✅ Elementos JSX Básicos
```jsx
<div>Hello, World!</div>
// → React.createElement("div", null, "Hello,World!")
```

### ✅ Elementos Self-Closing
```jsx
<img />
<br />
// → React.createElement("img", null)
// → React.createElement("br", null)
```

### ✅ Elementos Anidados
```jsx
<div>
  <span>Nested content</span>
</div>
// → React.createElement("div", null, React.createElement("span", null, "Nestedcontent"))
```

### ✅ Atributos JSX
```jsx
<div className="container" id="main">Content</div>
// → React.createElement("div", {"className": "container", "id": "main"}, "Content")
```

### ✅ Elementos Complejos
```jsx
<div className="main">
  <h1>Title</h1>
  <p>Paragraph</p>
</div>
// → React.createElement("div", {"className": "main"}, React.createElement("h1", null, "Title"), React.createElement("p", null, "Paragraph"))
```

## Uso

```go
package main

import (
    "github.com/xjslang/xjs/lexer"
    "github.com/xjslang/xjs/parser"
    jsxparser "github.com/xjslang/jsx-parser"
)

func main() {
    input := `let component = <div className="app">Hello, JSX!</div>`
    
    // Crear lexer y parser
    l := lexer.New(input)
    p := parser.New(l)
    
    // Registrar el middleware JSX
    p.UseExpressionHandler(jsxparser.ParseJsxExpression)
    
    // Parsear el programa
    ast := p.ParseProgram()
    
    // Obtener el código JavaScript transpilado
    output := ast.String()
    // → "let component = React.createElement("div", {"className": "app"}, "Hello,JSX!")"
}
```

## Ejemplos

Para ver ejemplos detallados de uso, consulta los **Example functions** incluidos en el código:

```bash
# Ver todos los ejemplos disponibles
go test -run Example -v

# Ejecutar un ejemplo específico
go test -run ExampleParseJsxExpression_nested -v
```

Los ejemplos cubren:
- **`ExampleParseJsxExpression`**: Uso básico
- **`ExampleParseJsxExpression_selfClosing`**: Elementos self-closing  
- **`ExampleParseJsxExpression_nested`**: Elementos anidados
- **`ExampleParseJsxExpression_withAttributes`**: JSX con atributos
- **`ExampleParseJsxExpression_complex`**: Casos complejos

## Arquitectura

El plugin utiliza el patrón middleware de XJS:

1. **Middleware Pattern**: Se registra como un `ExpressionHandler` que intercepta la parsing de expresiones
2. **Token Recognition**: Detecta la secuencia `<` + `IDENT` para iniciar el parsing JSX
3. **Recursive Parsing**: Maneja elementos anidados recursivamente
4. **Fallback**: Si no es JSX válido, pasa el control al parser siguiente en la cadena

### Componentes Principales

- **`JsxExpression`**: Representa un elemento JSX en el AST
- **`JsxAttribute`**: Representa un atributo de elemento JSX  
- **`JsxText`**: Representa contenido de texto dentro de elementos JSX
- **`ParseJsxExpression`**: Función middleware que maneja el parsing

## Limitaciones Actuales

- **Expresiones JavaScript**: No soporta `{expresion}` dentro de JSX
- **Fragmentos**: No soporta `<>...</>` (React Fragments)
- **Componentes**: Solo soporta elementos HTML nativos (tags en minúsculas)
- **Atributos complejos**: Solo soporta atributos con valores string literales

## Próximas Características

- [ ] Soporte para expresiones JavaScript `{variable}`
- [ ] React Fragments `<>...</>`
- [ ] Componentes React (PascalCase)
- [ ] Atributos con expresiones complejas
- [ ] Validación de nombres de tags HTML
- [ ] Mejor manejo de espacios en blanco en texto

## Testing

```bash
# Ejecutar todos los tests
go test -v

# Ejecutar solo los tests unitarios
go test -run TestJsxParser -v

# Ejecutar solo los ejemplos
go test -run Example -v

# Ver documentación con ejemplos
go doc .
```

El plugin incluye:
- **Tests unitarios comprehensivos** que cubren todos los casos de uso
- **Example functions** que sirven como documentación ejecutable y ejemplos de uso
