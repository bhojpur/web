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
	"fmt"
	"math"
)

type tagWidthratioNode struct {
	position     *Token
	current, max IEvaluator
	width        IEvaluator
	ctxName      string
}

func (node *tagWidthratioNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	current, err := node.current.Evaluate(ctx)
	if err != nil {
		return err
	}

	max, err := node.max.Evaluate(ctx)
	if err != nil {
		return err
	}

	width, err := node.width.Evaluate(ctx)
	if err != nil {
		return err
	}

	value := int(math.Ceil(current.Float()/max.Float()*width.Float() + 0.5))

	if node.ctxName == "" {
		writer.WriteString(fmt.Sprintf("%d", value))
	} else {
		ctx.Private[node.ctxName] = value
	}

	return nil
}

func tagWidthratioParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	widthratioNode := &tagWidthratioNode{
		position: start,
	}

	current, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	widthratioNode.current = current

	max, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	widthratioNode.max = max

	width, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	widthratioNode.width = width

	if arguments.MatchOne(TokenKeyword, "as") != nil {
		// Name follows
		nameToken := arguments.MatchType(TokenIdentifier)
		if nameToken == nil {
			return nil, arguments.Error("Expected name (identifier).", nil)
		}
		widthratioNode.ctxName = nameToken.Val
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed widthratio-tag arguments.", nil)
	}

	return widthratioNode, nil
}

func init() {
	RegisterTag("widthratio", tagWidthratioParser)
}
