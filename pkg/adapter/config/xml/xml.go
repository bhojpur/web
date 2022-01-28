// depend on github.com/bhojpur/web/pkg/x2j.
//
// go install github.com/bhojpur/web/pkg/x2j.
//
// Usage:
//  import(
//    _ "github.com/bhojpur/web/pkg/adapter/config/xml"
//      "github.com/bhojpur/web/pkg/config"
//  )
//
//  cnf, err := config.NewConfig("xml", "config.xml")
package xml

import (
	_ "github.com/bhojpur/web/pkg/core/config/xml"
)
