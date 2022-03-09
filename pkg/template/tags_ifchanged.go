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

type tagIfchangedNode struct {
	watchedExpr []IEvaluator
	lastValues  []*Value
	lastContent []byte
	thenWrapper *NodeWrapper
	elseWrapper *NodeWrapper
}

func (node *tagIfchangedNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	if len(node.watchedExpr) == 0 {
		// Check against own rendered body

		buf := bytes.NewBuffer(make([]byte, 0, 1024)) // 1 KiB
		err := node.thenWrapper.Execute(ctx, buf)
		if err != nil {
			return err
		}

		bufBytes := buf.Bytes()
		if !bytes.Equal(node.lastContent, bufBytes) {
			// Rendered content changed, output it
			writer.Write(bufBytes)
			node.lastContent = bufBytes
		}
	} else {
		nowValues := make([]*Value, 0, len(node.watchedExpr))
		for _, expr := range node.watchedExpr {
			val, err := expr.Evaluate(ctx)
			if err != nil {
				return err
			}
			nowValues = append(nowValues, val)
		}

		// Compare old to new values now
		changed := len(node.lastValues) == 0

		for idx, oldVal := range node.lastValues {
			if !oldVal.EqualValueTo(nowValues[idx]) {
				changed = true
				break // we can stop here because ONE value changed
			}
		}

		node.lastValues = nowValues

		if changed {
			// Render thenWrapper
			err := node.thenWrapper.Execute(ctx, writer)
			if err != nil {
				return err
			}
		} else {
			// Render elseWrapper
			err := node.elseWrapper.Execute(ctx, writer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func tagIfchangedParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	ifchangedNode := &tagIfchangedNode{}

	for arguments.Remaining() > 0 {
		// Parse condition
		expr, err := arguments.ParseExpression()
		if err != nil {
			return nil, err
		}
		ifchangedNode.watchedExpr = append(ifchangedNode.watchedExpr, expr)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Ifchanged-arguments are malformed.", nil)
	}

	// Wrap then/else-blocks
	wrapper, endargs, err := doc.WrapUntilTag("else", "endifchanged")
	if err != nil {
		return nil, err
	}
	ifchangedNode.thenWrapper = wrapper

	if endargs.Count() > 0 {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	if wrapper.Endtag == "else" {
		// if there's an else in the if-statement, we need the else-Block as well
		wrapper, endargs, err = doc.WrapUntilTag("endifchanged")
		if err != nil {
			return nil, err
		}
		ifchangedNode.elseWrapper = wrapper

		if endargs.Count() > 0 {
			return nil, endargs.Error("Arguments not allowed here.", nil)
		}
	}

	return ifchangedNode, nil
}

func init() {
	RegisterTag("ifchanged", tagIfchangedParser)
}
