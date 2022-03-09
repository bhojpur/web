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

type tagTemplateTagNode struct {
	content string
}

var templateTagMapping = map[string]string{
	"openblock":     "{%",
	"closeblock":    "%}",
	"openvariable":  "{{",
	"closevariable": "}}",
	"openbrace":     "{",
	"closebrace":    "}",
	"opencomment":   "{#",
	"closecomment":  "#}",
}

func (node *tagTemplateTagNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	writer.WriteString(node.content)
	return nil
}

func tagTemplateTagParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	ttNode := &tagTemplateTagNode{}

	if argToken := arguments.MatchType(TokenIdentifier); argToken != nil {
		output, found := templateTagMapping[argToken.Val]
		if !found {
			return nil, arguments.Error("Argument not found", argToken)
		}
		ttNode.content = output
	} else {
		return nil, arguments.Error("Identifier expected.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed templatetag-tag argument.", nil)
	}

	return ttNode, nil
}

func init() {
	RegisterTag("templatetag", tagTemplateTagParser)
}
