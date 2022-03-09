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

type tagSetNode struct {
	name       string
	expression IEvaluator
}

func (node *tagSetNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	// Evaluate expression
	value, err := node.expression.Evaluate(ctx)
	if err != nil {
		return err
	}

	ctx.Private[node.name] = value
	return nil
}

func tagSetParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	node := &tagSetNode{}

	// Parse variable name
	typeToken := arguments.MatchType(TokenIdentifier)
	if typeToken == nil {
		return nil, arguments.Error("Expected an identifier.", nil)
	}
	node.name = typeToken.Val

	if arguments.Match(TokenSymbol, "=") == nil {
		return nil, arguments.Error("Expected '='.", nil)
	}

	// Variable expression
	keyExpression, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	node.expression = keyExpression

	// Remaining arguments
	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed 'set'-tag arguments.", nil)
	}

	return node, nil
}

func init() {
	RegisterTag("set", tagSetParser)
}
