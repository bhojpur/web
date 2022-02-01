package swagger

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
	"github.com/bhojpur/web/pkg/swagger"
)

// Swagger list the resource
type Swagger swagger.Swagger

// Information Provides metadata about the API. The metadata can be used by the clients if needed.
type Information swagger.Information

// Contact information for the exposed API.
type Contact swagger.Contact

// License information for the exposed API.
type License swagger.License

// Item Describes the operations available on a single path.
type Item swagger.Item

// Operation Describes a single API operation on a path.
type Operation swagger.Operation

// Parameter Describes a single operation parameter.
type Parameter swagger.Parameter

// ParameterItems A limited subset of JSON-Schema's items object. It is used by parameter definitions that are not located in "body".
type ParameterItems swagger.ParameterItems

// Schema Object allows the definition of input and output data types.
type Schema swagger.Schema

// Propertie are taken from the JSON Schema definition but their definitions were adjusted to the Swagger Specification
type Propertie swagger.Propertie

// Response as they are returned from executing this operation.
type Response swagger.Response

// Security Allows the definition of a security scheme that can be used by the operations
type Security swagger.Security

// Tag Allows adding meta data to a single tag that is used by the Operation Object
type Tag swagger.Tag

// ExternalDocs include Additional external documentation
type ExternalDocs swagger.ExternalDocs
