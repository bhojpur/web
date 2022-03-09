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

type tagIfEqualNode struct {
	var1, var2  IEvaluator
	thenWrapper *NodeWrapper
	elseWrapper *NodeWrapper
}

func (node *tagIfEqualNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	r1, err := node.var1.Evaluate(ctx)
	if err != nil {
		return err
	}
	r2, err := node.var2.Evaluate(ctx)
	if err != nil {
		return err
	}

	result := r1.EqualValueTo(r2)

	if result {
		return node.thenWrapper.Execute(ctx, writer)
	}
	if node.elseWrapper != nil {
		return node.elseWrapper.Execute(ctx, writer)
	}
	return nil
}

func tagIfEqualParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	ifequalNode := &tagIfEqualNode{}

	// Parse two expressions
	var1, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	var2, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	ifequalNode.var1 = var1
	ifequalNode.var2 = var2

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("ifequal only takes 2 arguments.", nil)
	}

	// Wrap then/else-blocks
	wrapper, endargs, err := doc.WrapUntilTag("else", "endifequal")
	if err != nil {
		return nil, err
	}
	ifequalNode.thenWrapper = wrapper

	if endargs.Count() > 0 {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	if wrapper.Endtag == "else" {
		// if there's an else in the if-statement, we need the else-Block as well
		wrapper, endargs, err = doc.WrapUntilTag("endifequal")
		if err != nil {
			return nil, err
		}
		ifequalNode.elseWrapper = wrapper

		if endargs.Count() > 0 {
			return nil, endargs.Error("Arguments not allowed here.", nil)
		}
	}

	return ifequalNode, nil
}

func init() {
	RegisterTag("ifequal", tagIfEqualParser)
}
