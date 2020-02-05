package cherrypy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetJob(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "jobs_get", "success")

	res, err := c.Job(context.Background(), testSampleJobID)

	assert.NoError(t, err)
	assert.Equal(t, testSampleJobID, res.ID)
	assert.Equal(t, "cmd.run", res.Function)
	assert.Equal(t, "*", res.Target.(*ExpressionTarget).Expression)
	assert.Equal(t, Glob, res.Target.(*ExpressionTarget).Type)
	assert.Equal(t, "sudo_vagrant", res.User)
	assert.Equal(t, time.Date(2020, time.February, 2, 21, 2, 31, 414902000, time.UTC), res.StartTime)
	assert.Equal(t, "minion1", res.Minions[0])
	assert.Equal(t, "minion2", res.Minions[1])
	assert.Equal(t, 2, len(res.KWArguments))
	assert.Equal(t, "testy", res.KWArguments["test"])
	assert.Equal(t, "Can", res.KWArguments["complex_arg"].(map[string]interface{})["FIRST_NAME"])
	assert.Equal(t, "echo Hello", res.Arguments[0])
	assert.Equal(t, 1, len(res.Arguments))
}

func TestGetMissingJob(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "jobs_get", "missing")

	_, err := c.Job(context.Background(), "SampleMissingJobId")

	assert.Error(t, err)
}

func TestGetJobs(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "jobs_list", "success")

	res, err := c.Jobs(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 11, len(res))

	var job *Job
	for _, v := range res {
		if v.ID == testSampleJobID {
			job = &v
			break
		}
	}

	assert.NotNil(t, job)
	assert.Equal(t, testSampleJobID, job.ID)
	assert.Equal(t, "cmd.run", job.Function)
	assert.Equal(t, "*", job.Target.(*ExpressionTarget).Expression)
	assert.Equal(t, Glob, job.Target.(*ExpressionTarget).Type)
	assert.Equal(t, "sudo_vagrant", job.User)
	assert.Equal(t, time.Date(2020, time.February, 2, 21, 2, 31, 414902000, time.UTC), job.StartTime)
	assert.Equal(t, 2, len(job.KWArguments))
	assert.Equal(t, "testy", job.KWArguments["test"])
	assert.Equal(t, "Can", job.KWArguments["complex_arg"].(map[string]interface{})["FIRST_NAME"])
	assert.Equal(t, 1, len(job.Arguments))
	assert.Equal(t, "echo Hello", job.Arguments[0])
}
