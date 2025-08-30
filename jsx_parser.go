package jsxparser

import (
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// JsxExpression representa un elemento JSX
type JsxExpression struct {
	Token       token.Token      // el token '<'
	TagName     string           // nombre del tag (ej: "div", "span")
	Attributes  []JsxAttribute   // atributos del elemento
	Children    []ast.Expression // contenido hijo (texto o otros elementos JSX)
	SelfClosing bool             // true si es self-closing como <img />
}

// JsxAttribute representa un atributo de un elemento JSX
type JsxAttribute struct {
	Name  string         // nombre del atributo
	Value ast.Expression // valor del atributo (puede ser string o expresión)
}

// JsxText representa texto dentro de un elemento JSX
type JsxText struct {
	Token token.Token
	Value string
}

func (jt *JsxText) String() string {
	// Escapar comillas en el texto
	escaped := strings.ReplaceAll(jt.Value, "\"", "\\\"")
	return "\"" + escaped + "\""
}

func (jsx *JsxExpression) String() string {
	// Convertir JSX a JavaScript usando React.createElement
	if jsx.SelfClosing || len(jsx.Children) == 0 {
		// Elemento sin hijos: React.createElement("tagName", props)
		return jsx.toCreateElement()
	} else {
		// Elemento con hijos: React.createElement("tagName", props, ...children)
		return jsx.toCreateElementWithChildren()
	}
}

func (jsx *JsxExpression) toCreateElement() string {
	props := jsx.attributesToProps()
	if props == "null" {
		return "React.createElement(\"" + jsx.TagName + "\", null)"
	}
	return "React.createElement(\"" + jsx.TagName + "\", " + props + ")"
}

func (jsx *JsxExpression) toCreateElementWithChildren() string {
	props := jsx.attributesToProps()
	children := jsx.childrenToString()

	result := "React.createElement(\"" + jsx.TagName + "\", " + props
	if children != "" {
		result += ", " + children
	}
	result += ")"
	return result
}

func (jsx *JsxExpression) attributesToProps() string {
	if len(jsx.Attributes) == 0 {
		return "null"
	}

	var props []string
	for _, attr := range jsx.Attributes {
		props = append(props, "\""+attr.Name+"\": "+attr.Value.String())
	}
	return "{" + strings.Join(props, ", ") + "}"
}

func (jsx *JsxExpression) childrenToString() string {
	if len(jsx.Children) == 0 {
		return ""
	}

	var children []string
	for _, child := range jsx.Children {
		children = append(children, child.String())
	}
	return strings.Join(children, ", ")
}

func ParseJsxExpression(p *parser.Parser, precedence int, next func(*parser.Parser, int) ast.Expression) ast.Expression {
	// Solo procesar si encontramos '<' seguido de un identificador
	if p.CurrentToken.Type != token.LT || p.PeekToken.Type != token.IDENT {
		return next(p, precedence)
	}

	jsx := &JsxExpression{
		Token:      p.CurrentToken,
		Attributes: []JsxAttribute{},
		Children:   []ast.Expression{},
	}

	// Consumir '<'
	p.NextToken()

	// Obtener el nombre del tag
	jsx.TagName = p.CurrentToken.Literal
	p.NextToken()

	// Procesar atributos
	for p.CurrentToken.Type == token.IDENT {
		attr := JsxAttribute{
			Name: p.CurrentToken.Literal,
		}
		p.NextToken()

		// Verificar si tiene valor (atributo="valor")
		if p.CurrentToken.Type == token.ASSIGN {
			p.NextToken()
			if p.CurrentToken.Type == token.STRING {
				// Valor de atributo como string literal
				attr.Value = &ast.StringLiteral{
					Token: p.CurrentToken,
					Value: p.CurrentToken.Literal,
				}
				p.NextToken()
			} else {
				// Por simplicidad, tratamos otros valores como strings
				attr.Value = &ast.StringLiteral{
					Token: p.CurrentToken,
					Value: p.CurrentToken.Literal,
				}
				p.NextToken()
			}
		} else {
			// Atributo booleano sin valor (ej: disabled)
			attr.Value = &ast.BooleanLiteral{
				Token: token.Token{Type: token.TRUE, Literal: "true"},
				Value: true,
			}
		}

		jsx.Attributes = append(jsx.Attributes, attr)
	}

	// Verificar si es self-closing
	if p.CurrentToken.Type == token.DIVIDE && p.PeekToken.Type == token.GT {
		jsx.SelfClosing = true
		p.NextToken() // consume '/'
		p.NextToken() // consume '>'
		return jsx
	}

	// Consumir '>' de apertura
	if p.CurrentToken.Type != token.GT {
		p.AddError("expected '>' after tag name")
		return nil
	}
	p.NextToken()

	// Procesar contenido hasta encontrar el tag de cierre
	var textBuffer []string

	for p.CurrentToken.Type != token.EOF {
		// Verificar si es el inicio de un tag de cierre
		if p.CurrentToken.Type == token.LT && p.PeekToken.Type == token.DIVIDE {
			// Si tenemos texto acumulado, agregarlo como un nodo de texto
			if len(textBuffer) > 0 {
				text := &JsxText{
					Token: p.CurrentToken,
					Value: strings.Join(textBuffer, ""),
				}
				jsx.Children = append(jsx.Children, text)
			}

			// Consumir '</'
			p.NextToken()
			p.NextToken()

			// Verificar que el nombre del tag de cierre coincida
			if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == jsx.TagName {
				p.NextToken() // consume tag name
				if p.CurrentToken.Type == token.GT {
					p.NextToken() // consume '>'
					break
				}
			}
			p.AddError("malformed closing tag")
			return nil
		}

		// Si encontramos otro elemento JSX anidado
		if p.CurrentToken.Type == token.LT && p.PeekToken.Type == token.IDENT {
			// Si tenemos texto acumulado, agregarlo antes del elemento JSX
			if len(textBuffer) > 0 {
				text := &JsxText{
					Token: p.CurrentToken,
					Value: strings.Join(textBuffer, ""),
				}
				jsx.Children = append(jsx.Children, text)
				textBuffer = textBuffer[:0] // clear buffer
			}

			childJsx := ParseJsxExpression(p, precedence, next)
			if childJsx != nil {
				jsx.Children = append(jsx.Children, childJsx)
			}
			continue
		}

		// Acumular texto (identificadores, strings, signos de puntuación, espacios)
		if p.CurrentToken.Literal != "" &&
			p.CurrentToken.Type != token.LT &&
			p.CurrentToken.Type != token.GT {
			textBuffer = append(textBuffer, p.CurrentToken.Literal)
		}
		p.NextToken()
	}

	return jsx
}
