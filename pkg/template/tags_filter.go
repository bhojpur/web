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
	"bytes"
)

type nodeFilterCall struct {
	name      string
	paramExpr IEvaluator
}

type tagFilterNode struct {
	position    *Token
	bodyWrapper *NodeWrapper
	filterChain []*nodeFilterCall
}

func (node *tagFilterNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	temp := bytes.NewBuffer(make([]byte, 0, 1024)) // 1 KiB size

	err := node.bodyWrapper.Execute(ctx, temp)
	if err != nil {
		return err
	}

	value := AsValue(temp.String())

	for _, call := range node.filterChain {
		var param *Value
		if call.paramExpr != nil {
			param, err = call.paramExpr.Evaluate(ctx)
			if err != nil {
				return err
			}
		} else {
			param = AsValue(nil)
		}
		value, err = ApplyFilter(call.name, value, param)
		if err != nil {
			return ctx.Error(err.Error(), node.position)
		}
	}

	writer.WriteString(value.String())

	return nil
}

func tagFilterParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	filterNode := &tagFilterNode{
		position: start,
	}

	wrapper, _, err := doc.WrapUntilTag("endfilter")
	if err != nil {
		return nil, err
	}
	filterNode.bodyWrapper = wrapper

	for arguments.Remaining() > 0 {
		filterCall := &nodeFilterCall{}

		nameToken := arguments.MatchType(TokenIdentifier)
		if nameToken == nil {
			return nil, arguments.Error("Expected a filter name (identifier).", nil)
		}
		filterCall.name = nameToken.Val

		if arguments.MatchOne(TokenSymbol, ":") != nil {
			// Filter parameter
			// NOTICE: we can't use ParseExpression() here, because it would parse the next filter "|..." as well in the argument list
			expr, err := arguments.parseVariableOrLiteral()
			if err != nil {
				return nil, err
			}
			filterCall.paramExpr = expr
		}

		filterNode.filterChain = append(filterNode.filterChain, filterCall)

		if arguments.MatchOne(TokenSymbol, "|") == nil {
			break
		}
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed filter-tag arguments.", nil)
	}

	return filterNode, nil
}

func init() {
	RegisterTag("filter", tagFilterParser)
}
