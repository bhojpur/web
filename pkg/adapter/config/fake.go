package config

import (
	"github.com/bhojpur/web/pkg/core/config"
)

// NewFakeConfig return a fake Configer
func NewFakeConfig() Configer {
	new := config.NewFakeConfig()
	return &newToOldConfigerAdapter{delegate: new}
}
