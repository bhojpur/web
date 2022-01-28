// depend on github.com/bhojpur/web/pkg/core/config/yaml
//
// go install github.com/bhojpur/web
//
// Usage:
//  import(
//   _ "github.com/bhojpur/web/pkg/adapter/config/yaml"
//     "github.com/bhojpur/web/pkg/config"
//  )
//
//  cnf, err := config.NewConfig("yaml", "config.yaml")
package yaml

import (
	_ "github.com/bhojpur/web/pkg/core/config/yaml"
)
