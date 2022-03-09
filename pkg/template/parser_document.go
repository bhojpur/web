package template

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Doc = { ( Filter | Tag | HTML ) }
func (p *Parser) parseDocElement() (INode, *Error) {
	t := p.Current()

	switch t.Typ {
	case TokenHTML:
		n := &nodeHTML{token: t}
		left := p.PeekTypeN(-1, TokenSymbol)
		right := p.PeekTypeN(1, TokenSymbol)
		n.trimLeft = left != nil && left.TrimWhitespaces
		n.trimRight = right != nil && right.TrimWhitespaces
		p.Consume() // consume HTML element
		return n, nil
	case TokenSymbol:
		switch t.Val {
		case "{{":
			// parse variable
			variable, err := p.parseVariableElement()
			if err != nil {
				return nil, err
			}
			return variable, nil
		case "{%":
			// parse tag
			tag, err := p.parseTagElement()
			if err != nil {
				return nil, err
			}
			return tag, nil
		}
	}
	return nil, p.Error("Unexpected token (only HTML/tags/filters in templates allowed)", t)
}

func (tpl *Template) parse() *Error {
	tpl.parser = newParser(tpl.name, tpl.tokens, tpl)
	doc, err := tpl.parser.parseDocument()
	if err != nil {
		return err
	}
	tpl.root = doc
	return nil
}

func (p *Parser) parseDocument() (*nodeDocument, *Error) {
	doc := &nodeDocument{}

	for p.Remaining() > 0 {
		node, err := p.parseDocElement()
		if err != nil {
			return nil, err
		}
		doc.Nodes = append(doc.Nodes, node)
	}

	return doc, nil
}
