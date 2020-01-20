package cherrypy

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	testSampleJobID = "20200120025629463543"
)

func TestGetJob(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/jobs/"+testSampleJobID, "job_get_success")

	res, err := c.Job(testSampleJobID)

	assert.NoError(t, err)
	assert.Equal(t, testSampleJobID, res.ID)
	assert.Equal(t, "cmd.run", res.Function)
	assert.Equal(t, "*", res.Target.Target())
	assert.Equal(t, "glob", res.Target.Type())
	assert.Equal(t, "test_user", res.User)
	assert.Equal(t, time.Date(2020, time.January, 20, 2, 56, 29, 463543000, time.UTC), res.StartTime)
	assert.Equal(t, "minion1", res.Minions[0])
	assert.Equal(t, "minion2", res.Minions[1])
	assert.Equal(t, 2, len(res.KWArguments))
	assert.Equal(t, "testy", res.KWArguments["test"])
	assert.Equal(t, "Can", res.KWArguments["complex_arg"].(map[string]interface{})["FIRST_NAME"])
	assert.Equal(t, "echo Hello", res.Arguments[0])
	assert.Equal(t, 1, len(res.Arguments))
}

func TestGetMissingJob(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/jobs/SampleMissingJobId", "job_get_missing")

	_, err := c.Job("SampleMissingJobId")

	assert.Error(t, err)
}

func TestGetJobs(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/jobs", "job_list_success")

	res, err := c.Jobs()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	job := res[1]
	assert.NotNil(t, job)
	assert.Equal(t, testSampleJobID, job.ID)
	assert.Equal(t, "cmd.run", job.Function)
	assert.Equal(t, "*", job.Target.Target())
	assert.Equal(t, "glob", job.Target.Type())
	assert.Equal(t, "test_user", job.User)
	assert.Equal(t, time.Date(2020, time.January, 20, 2, 56, 29, 463543000, time.UTC), job.StartTime)
	assert.Equal(t, 2, len(job.KWArguments))
	assert.Equal(t, "testy", job.KWArguments["test"])
	assert.Equal(t, "Can", job.KWArguments["complex_arg"].(map[string]interface{})["FIRST_NAME"])
	assert.Equal(t, 1, len(job.Arguments))
	assert.Equal(t, "echo Hello", job.Arguments[0])
}
