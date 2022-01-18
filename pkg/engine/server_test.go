package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpServerWithCfg(t *testing.T) {

	BasConfig.AppName = "Before"
	svr := NewHttpServerWithCfg(BasConfig)
	svr.Cfg.AppName = "hello"
	assert.Equal(t, "hello", BasConfig.AppName)

}
