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

type tagAutoescapeNode struct {
	wrapper    *NodeWrapper
	autoescape bool
}

func (node *tagAutoescapeNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	old := ctx.Autoescape
	ctx.Autoescape = node.autoescape

	err := node.wrapper.Execute(ctx, writer)
	if err != nil {
		return err
	}

	ctx.Autoescape = old

	return nil
}

func tagAutoescapeParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	autoescapeNode := &tagAutoescapeNode{}

	wrapper, _, err := doc.WrapUntilTag("endautoescape")
	if err != nil {
		return nil, err
	}
	autoescapeNode.wrapper = wrapper

	modeToken := arguments.MatchType(TokenIdentifier)
	if modeToken == nil {
		return nil, arguments.Error("A mode is required for autoescape-tag.", nil)
	}
	if modeToken.Val == "on" {
		autoescapeNode.autoescape = true
	} else if modeToken.Val == "off" {
		autoescapeNode.autoescape = false
	} else {
		return nil, arguments.Error("Only 'on' or 'off' is valid as an autoescape-mode.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed autoescape-tag arguments.", nil)
	}

	return autoescapeNode, nil
}

func init() {
	RegisterTag("autoescape", tagAutoescapeParser)
}
