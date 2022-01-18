package bean

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeTypeAdapter_DefaultValue(t *testing.T) {
	typeAdapter := &TimeTypeAdapter{Layout: "2018-03-26 15:04:05"}
	tm, err := typeAdapter.DefaultValue(context.Background(), "2018-02-26 12:34:11")
	assert.Nil(t, err)
	assert.NotNil(t, tm)
}
