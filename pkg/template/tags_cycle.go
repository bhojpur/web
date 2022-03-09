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

type tagCycleValue struct {
	node  *tagCycleNode
	value *Value
}

type tagCycleNode struct {
	position *Token
	args     []IEvaluator
	idx      int
	asName   string
	silent   bool
}

func (cv *tagCycleValue) String() string {
	return cv.value.String()
}

func (node *tagCycleNode) Execute(ctx *ExecutionContext, writer TemplateWriter) *Error {
	item := node.args[node.idx%len(node.args)]
	node.idx++

	val, err := item.Evaluate(ctx)
	if err != nil {
		return err
	}

	if t, ok := val.Interface().(*tagCycleValue); ok {
		// {% cycle "test1" "test2"
		// {% cycle cycleitem %}

		// Update the cycle value with next value
		item := t.node.args[t.node.idx%len(t.node.args)]
		t.node.idx++

		val, err := item.Evaluate(ctx)
		if err != nil {
			return err
		}

		t.value = val

		if !t.node.silent {
			writer.WriteString(val.String())
		}
	} else {
		// Regular call

		cycleValue := &tagCycleValue{
			node:  node,
			value: val,
		}

		if node.asName != "" {
			ctx.Private[node.asName] = cycleValue
		}
		if !node.silent {
			writer.WriteString(val.String())
		}
	}

	return nil
}

// HINT: We're not supporting the old comma-separated list of expressions argument-style
func tagCycleParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	cycleNode := &tagCycleNode{
		position: start,
	}

	for arguments.Remaining() > 0 {
		node, err := arguments.ParseExpression()
		if err != nil {
			return nil, err
		}
		cycleNode.args = append(cycleNode.args, node)

		if arguments.MatchOne(TokenKeyword, "as") != nil {
			// as

			nameToken := arguments.MatchType(TokenIdentifier)
			if nameToken == nil {
				return nil, arguments.Error("Name (identifier) expected after 'as'.", nil)
			}
			cycleNode.asName = nameToken.Val

			if arguments.MatchOne(TokenIdentifier, "silent") != nil {
				cycleNode.silent = true
			}

			// Now we're finished
			break
		}
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed cycle-tag.", nil)
	}

	return cycleNode, nil
}

func init() {
	RegisterTag("cycle", tagCycleParser)
}
