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

import (
	"time"
)

type tagNowNode struct {
	position *Token
	format   string
	fake     bool
}

func (node *tagNowNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	var t time.Time
	if node.fake {
		t = time.Date(2014, time.February, 05, 18, 31, 45, 00, time.UTC)
	} else {
		t = time.Now()
	}

	writer.WriteString(t.Format(node.format))

	return nil
}

func tagNowParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	nowNode := &tagNowNode{
		position: start,
	}

	formatToken := arguments.MatchType(TokenString)
	if formatToken == nil {
		return nil, arguments.Error("Expected a format string.", nil)
	}
	nowNode.format = formatToken.Val

	if arguments.MatchOne(TokenIdentifier, "fake") != nil {
		nowNode.fake = true
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed now-tag arguments.", nil)
	}

	return nowNode, nil
}

func init() {
	RegisterTag("now", tagNowParser)
}
