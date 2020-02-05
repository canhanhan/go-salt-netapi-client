package cherrypy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsSuccess(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "stats", "success")

	res, err := c.Stats(context.Background())

	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
