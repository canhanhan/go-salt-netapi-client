package cherrypy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSingleOnlineMinion(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions/minion1", "minions_get_success")

	res, err := c.Minion("minion1")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Grains)
}

func TestGetSingleOfflineMinion(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions/minion2", "minions_get_offline")

	res, err := c.Minion("minion2")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, res.Grains)
}

func TestGetSingleMissingMinion(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions/minion3", "minions_get_missing")

	res, err := c.Minion("minion3")

	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestSubmitSingleJob(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions", "minions_submit_single")

	res, err := c.SubmitJob(MinionJob{
		Target:   "minion1",
		Function: "test.ping",
	})

	assert.NoError(t, err)
	assert.Contains(t, res.Minions, "minion1")
	assert.NotEmpty(t, res.JobID)
}

func TestSubmitSingleJobToOfflineMinion(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions", "minions_submit_offline")

	res, err := c.SubmitJob(MinionJob{
		Target:   "minion2",
		Function: "test.ping",
	})

	assert.NoError(t, err)
	assert.Contains(t, res.Minions, "minion2")
	assert.NotEmpty(t, res.JobID)
}

func TestSubmitSingleJobToMissingMinion(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions", "minions_submit_missing")

	res, err := c.SubmitJob(MinionJob{
		Target:   "minion3",
		Function: "test.ping",
	})

	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestSubmitMultipleJobs(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions", "minions_submit_multiple")

	res, err := c.SubmitJobs([]MinionJob{
		MinionJob{
			Target:   "minion1",
			Function: "test.ping",
		},
		MinionJob{
			Target:   "minion1",
			Function: "test.ping",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Contains(t, res[0].Minions, "minion1")
	assert.NotEmpty(t, res[0].JobID)
	assert.Contains(t, res[1].Minions, "minion1")
	assert.NotEmpty(t, res[1].JobID)
}

func TestSubmitMultipleJobToOfflineMinion(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions", "minions_submit_multiple_offline")

	res, err := c.SubmitJobs([]MinionJob{
		MinionJob{
			Target:   "minion2",
			Function: "test.ping",
		},
		MinionJob{
			Target:   "minion2",
			Function: "test.ping",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Contains(t, res[0].Minions, "minion2")
	assert.NotEmpty(t, res[0].JobID)
	assert.Contains(t, res[1].Minions, "minion2")
	assert.NotEmpty(t, res[1].JobID)
}

func TestSubmitMuiltipleJobsToMissingMinion(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions", "minions_submit_multiple_missing")

	res, err := c.SubmitJobs([]MinionJob{
		MinionJob{
			Target:   "minion3",
			Function: "test.ping",
		},
		MinionJob{
			Target:   "minion3",
			Function: "test.ping",
		},
	})

	assert.NoError(t, err)
	assert.Empty(t, res)
}

func TestSubmitMuiltipleJobsToMixedStatusMinions(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/minions", "minions_submit_multiple_mixed")

	res, err := c.SubmitJobs([]MinionJob{
		MinionJob{
			Target:   "minion1",
			Function: "test.ping",
		},
		MinionJob{
			Target:   "minion2",
			Function: "test.ping",
		},
		MinionJob{
			Target:   "minion3",
			Function: "test.ping",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Contains(t, res[0].Minions, "minion1")
	assert.NotEmpty(t, res[0].JobID)
	assert.Contains(t, res[1].Minions, "minion2")
	assert.NotEmpty(t, res[1].JobID)
}
