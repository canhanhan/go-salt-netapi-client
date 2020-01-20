package cherrypy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHookSuccess(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/hook/test", "hook_success")

	err := c.Hook("test", nil)

	assert.NoError(t, err)
}

func TestHookFailure(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/hook/test", "hook_failure")

	err := c.Hook("test", nil)

	assert.Error(t, err)
}
