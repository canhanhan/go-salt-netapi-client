package cherrypy

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestStatsSuccess(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()	
	handleJSONRequest(mux, "/stats", "stats_success")
	
	res, err := c.Stats()

	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
