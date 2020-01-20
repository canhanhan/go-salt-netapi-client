package cherrypy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunLocalCommand(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/run", "run_local_success")

	cmd := Command{
		Client:   "local",
		Target:   "minion1",
		Function: "test.ping",
	}

	res, err := c.RunJob(cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestRunLocalCommands(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/run", "run_local_multiple_success")

	cmds := []Command{
		Command{
			Client:   "local",
			Target:   "minion1",
			Function: "test.ping",
		},
		Command{
			Client:   "local",
			Target:   "minion1",
			Function: "test.ping",
		},
	}

	res, err := c.RunJobs(cmds)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}

func TestRunWheelCommand(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/run", "run_wheel_success")

	cmd := Command{
		Client:   "wheel",
		Function: "minions.connected",
	}

	res, err := c.RunJob(cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}
