package jsxparser

import (
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// JsxExpression represents a JSX element
type JsxExpression struct {
	Token       token.Token      // the '<' token
	TagName     string           // tag name (e.g., "div", "span")
	Attributes  []JsxAttribute   // element attributes
	Children    []ast.Expression // child content (text or other JSX elements)
	SelfClosing bool             // true if self-closing like <img />
}

// JsxAttribute represents an attribute of a JSX element
type JsxAttribute struct {
	Name  string         // attribute name
	Value ast.Expression // attribute value (can be string or expression)
}

// JsxText represents text inside a JSX element
type JsxText struct {
	Token token.Token
	Value string
}

func (jt *JsxText) WriteTo(b *strings.Builder) {
	// Escape quotes in the text
	b.WriteRune('"')
	for _, char := range jt.Value {
		if char == '"' {
			b.WriteString("\\\"")
		} else {
			b.WriteRune(char)
		}
	}
	b.WriteRune('"')
}

func (jsx *JsxExpression) WriteTo(b *strings.Builder) {
	// Convert JSX to JavaScript using React.createElement
	if jsx.SelfClosing || len(jsx.Children) == 0 {
		// Element without children: React.createElement("tagName", props)
		jsx.writeCreateElement(b)
	} else {
		// Element with children: React.createElement("tagName", props, ...children)
		jsx.writeCreateElementWithChildren(b)
	}
}

func (jsx *JsxExpression) writeCreateElement(b *strings.Builder) {
	b.WriteString("React.createElement(\"")
	b.WriteString(jsx.TagName)
	b.WriteString("\", ")
	jsx.writeAttributesToProps(b)
	b.WriteRune(')')
}

func (jsx *JsxExpression) writeCreateElementWithChildren(b *strings.Builder) {
	b.WriteString("React.createElement(\"")
	b.WriteString(jsx.TagName)
	b.WriteString("\", ")
	jsx.writeAttributesToProps(b)

	if len(jsx.Children) > 0 {
		b.WriteString(", ")
		jsx.writeChildrenToString(b)
	}
	b.WriteRune(')')
}

func (jsx *JsxExpression) writeAttributesToProps(b *strings.Builder) {
	if len(jsx.Attributes) == 0 {
		b.WriteString("null")
		return
	}

	b.WriteRune('{')
	for i, attr := range jsx.Attributes {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteRune('"')
		b.WriteString(attr.Name)
		b.WriteString("\": ")
		attr.Value.WriteTo(b)
	}
	b.WriteRune('}')
}

func (jsx *JsxExpression) writeChildrenToString(b *strings.Builder) {
	for i, child := range jsx.Children {
		if i > 0 {
			b.WriteString(", ")
		}
		child.WriteTo(b)
	}
}

func ParseJsxExpression(p *parser.Parser, next func(left ast.Expression) ast.Expression) ast.Expression {
	// Only process if we find '<' followed by an identifier
	if p.CurrentToken.Type != token.LT || p.PeekToken.Type != token.IDENT {
		return next(nil)
	}

	jsx := &JsxExpression{
		Token:      p.CurrentToken,
		Attributes: []JsxAttribute{},
		Children:   []ast.Expression{},
	}

	// Consume '<'
	p.NextToken()

	// Get the tag name
	jsx.TagName = p.CurrentToken.Literal
	p.NextToken()

	// Process attributes
	for p.CurrentToken.Type == token.IDENT {
		attr := JsxAttribute{
			Name: p.CurrentToken.Literal,
		}
		p.NextToken()

		// Check if it has a value (attribute="value")
		if p.CurrentToken.Type == token.ASSIGN {
			p.NextToken()
			if p.CurrentToken.Type == token.STRING {
				// Attribute value as string literal
				attr.Value = &ast.StringLiteral{
					Token: p.CurrentToken,
					Value: p.CurrentToken.Literal,
				}
				p.NextToken()
			} else {
				// For simplicity, treat other values as strings
				attr.Value = &ast.StringLiteral{
					Token: p.CurrentToken,
					Value: p.CurrentToken.Literal,
				}
				p.NextToken()
			}
		} else {
			// Boolean attribute without value (e.g., disabled)
			attr.Value = &ast.BooleanLiteral{
				Token: token.Token{Type: token.TRUE, Literal: "true"},
				Value: true,
			}
		}

		jsx.Attributes = append(jsx.Attributes, attr)
	}

	// Check if self-closing
	if p.CurrentToken.Type == token.DIVIDE && p.PeekToken.Type == token.GT {
		jsx.SelfClosing = true
		p.NextToken() // consume '/'
		p.NextToken() // consume '>'
		return jsx
	}

	// Consume opening '>'
	if p.CurrentToken.Type != token.GT {
		p.AddError("expected '>' after tag name")
		return nil
	}
	p.NextToken()

	// Process content until finding the closing tag
	var textBuffer []string

	for p.CurrentToken.Type != token.EOF {
		// Check if it's the start of a closing tag
		if p.CurrentToken.Type == token.LT && p.PeekToken.Type == token.DIVIDE {
			// If we have accumulated text, add it as a text node
			if len(textBuffer) > 0 {
				text := &JsxText{
					Token: p.CurrentToken,
					Value: strings.Join(textBuffer, ""),
				}
				jsx.Children = append(jsx.Children, text)
			}

			// Consume '</'
			p.NextToken()
			p.NextToken()

			// Check that the closing tag name matches
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

		// If we find another nested JSX element
		if p.CurrentToken.Type == token.LT && p.PeekToken.Type == token.IDENT {
			// If we have accumulated text, add it before the JSX element
			if len(textBuffer) > 0 {
				text := &JsxText{
					Token: p.CurrentToken,
					Value: strings.Join(textBuffer, ""),
				}
				jsx.Children = append(jsx.Children, text)
				textBuffer = textBuffer[:0] // clear buffer
			}

			childJsx := ParseJsxExpression(p, next)
			if childJsx != nil {
				jsx.Children = append(jsx.Children, childJsx)
			}
			continue
		}

		// Accumulate text (identifiers, strings, punctuation, spaces)
		if p.CurrentToken.Literal != "" &&
			p.CurrentToken.Type != token.LT &&
			p.CurrentToken.Type != token.GT {
			textBuffer = append(textBuffer, p.CurrentToken.Literal)
		}
		p.NextToken()
	}

	return jsx
}
