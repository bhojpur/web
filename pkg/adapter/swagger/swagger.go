package swagger

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
