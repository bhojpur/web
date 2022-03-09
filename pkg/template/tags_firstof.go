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

type tagFirstofNode struct {
	position *Token
	args     []IEvaluator
}

func (node *tagFirstofNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	for _, arg := range node.args {
		val, err := arg.Evaluate(ctx)
		if err != nil {
			return err
		}

		if val.IsTrue() {
			if ctx.Autoescape && !arg.FilterApplied("safe") {
				val, err = ApplyFilter("escape", val, nil)
				if err != nil {
					return err
				}
			}

			writer.WriteString(val.String())
			return nil
		}
	}

	return nil
}

func tagFirstofParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	firstofNode := &tagFirstofNode{
		position: start,
	}

	for arguments.Remaining() > 0 {
		node, err := arguments.ParseExpression()
		if err != nil {
			return nil, err
		}
		firstofNode.args = append(firstofNode.args, node)
	}

	return firstofNode, nil
}

func init() {
	RegisterTag("firstof", tagFirstofParser)
}
