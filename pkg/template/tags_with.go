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

type tagWithNode struct {
	withPairs map[string]IEvaluator
	wrapper   *NodeWrapper
}

func (node *tagWithNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	//new context for block
	withctx := NewChildExecutionContext(ctx)

	// Put all custom with-pairs into the context
	for key, value := range node.withPairs {
		val, err := value.Evaluate(ctx)
		if err != nil {
			return err
		}
		withctx.Private[key] = val
	}

	return node.wrapper.Execute(withctx, writer)
}

func tagWithParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	withNode := &tagWithNode{
		withPairs: make(map[string]IEvaluator),
	}

	if arguments.Count() == 0 {
		return nil, arguments.Error("Tag 'with' requires at least one argument.", nil)
	}

	wrapper, endargs, err := doc.WrapUntilTag("endwith")
	if err != nil {
		return nil, err
	}
	withNode.wrapper = wrapper

	if endargs.Count() > 0 {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	// Scan through all arguments to see which style the user uses (old or new style).
	// If we find any "as" keyword we will enforce old style; otherwise we will use new style.
	oldStyle := false // by default we're using the new_style
	for i := 0; i < arguments.Count(); i++ {
		if arguments.PeekN(i, TokenKeyword, "as") != nil {
			oldStyle = true
			break
		}
	}

	for arguments.Remaining() > 0 {
		if oldStyle {
			valueExpr, err := arguments.ParseExpression()
			if err != nil {
				return nil, err
			}
			if arguments.Match(TokenKeyword, "as") == nil {
				return nil, arguments.Error("Expected 'as' keyword.", nil)
			}
			keyToken := arguments.MatchType(TokenIdentifier)
			if keyToken == nil {
				return nil, arguments.Error("Expected an identifier", nil)
			}
			withNode.withPairs[keyToken.Val] = valueExpr
		} else {
			keyToken := arguments.MatchType(TokenIdentifier)
			if keyToken == nil {
				return nil, arguments.Error("Expected an identifier", nil)
			}
			if arguments.Match(TokenSymbol, "=") == nil {
				return nil, arguments.Error("Expected '='.", nil)
			}
			valueExpr, err := arguments.ParseExpression()
			if err != nil {
				return nil, err
			}
			withNode.withPairs[keyToken.Val] = valueExpr
		}
	}

	return withNode, nil
}

func init() {
	RegisterTag("with", tagWithParser)
}
