package berror

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
)

// pre define code

// Unknown indicates got some error which is not defined
var Unknown = DefineCode(5000001, "error", "Unknown", fmt.Sprintf(`
Unknown error code. Usually you will see this code in three cases:
1. You forget to define Code or function DefineCode not being executed;
2. This is not Bhojpur Web's error but you call FromError();
3. Bhojpur Web got unexpected error and don't know how to handle it, and then return Unknown error

A common practice to DefineCode looks like:
%s

In this way, you may forget to import this package, and got Unknown error. 

Sometimes, you believe you got Bhojpur Web error, but actually you don't, and then you call FromError(err)

`, goCodeBlock(`
import your_package

func init() {
    DefineCode(5100100, "your_module", "detail")
    // ...
}
`)))

func goCodeBlock(code string) string {
	return codeBlock("go", code)
}

func codeBlock(lan string, code string) string {
	return fmt.Sprintf("```%s\n%s\n```", lan, code)
}
