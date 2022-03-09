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

type tagExtendsNode struct {
	filename string
}

func (node *tagExtendsNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	return nil
}

func tagExtendsParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	extendsNode := &tagExtendsNode{}

	if doc.template.level > 1 {
		return nil, arguments.Error("The 'extends' tag can only defined on root level.", start)
	}

	if doc.template.parent != nil {
		// Already one parent
		return nil, arguments.Error("This template has already one parent.", start)
	}

	if filenameToken := arguments.MatchType(TokenString); filenameToken != nil {
		// prepared, static template

		// Get parent's filename
		parentFilename := doc.template.set.resolveFilename(doc.template, filenameToken.Val)

		// Parse the parent
		parentTemplate, err := doc.template.set.FromFile(parentFilename)
		if err != nil {
			return nil, err.(*Error)
		}

		// Keep track of things
		parentTemplate.child = doc.template
		doc.template.parent = parentTemplate
		extendsNode.filename = parentFilename
	} else {
		return nil, arguments.Error("Tag 'extends' requires a template filename as string.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Tag 'extends' does only take 1 argument.", nil)
	}

	return extendsNode, nil
}

func init() {
	RegisterTag("extends", tagExtendsParser)
}
