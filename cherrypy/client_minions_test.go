package cherrypy

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSingleOnlineMinion(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_get", "success")

	res, err := c.Minion(context.Background(), "minion1")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Grains)
}

func TestGetSingleOfflineMinion(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_get", "offline")

	res, err := c.Minion(context.Background(), "minion2")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, res.Grains)
}

func TestGetSingleMissingMinion(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_get", "missing")

	res, err := c.Minion(context.Background(), "minion3")
	if !errors.Is(err, ErrorMinionNotFound) {
		t.Fatal(err)
	}

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestSubmitSingleJob(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_submit", "single")

	res, err := c.SubmitJob(context.Background(), MinionJob{
		Target:   ExpressionTarget{Expression: "minion1", Type: Glob},
		Function: "test.ping",
	})

	assert.NoError(t, err)
	assert.Contains(t, res.Minions, "minion1")
	assert.NotEmpty(t, res.ID)
}

func TestSubmitSingleJobToOfflineMinion(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_submit", "offline")

	res, err := c.SubmitJob(context.Background(), MinionJob{
		Target:   ExpressionTarget{Expression: "minion2", Type: Glob},
		Function: "test.ping",
	})

	assert.NoError(t, err)
	assert.Contains(t, res.Minions, "minion2")
	assert.NotEmpty(t, res.ID)
}

func TestSubmitSingleJobToMissingMinion(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_submit", "missing")

	res, err := c.SubmitJob(context.Background(), MinionJob{
		Target:   ExpressionTarget{Expression: "minion3", Type: Glob},
		Function: "test.ping",
	})

	assert.NoError(t, err)
	assert.Empty(t, res.ID)
}

func TestSubmitMultipleJobs(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_submit", "multiple")

	res, err := c.SubmitJobs(context.Background(), []MinionJob{
		MinionJob{
			Target:   ExpressionTarget{Expression: "minion1", Type: Glob},
			Function: "test.ping",
		},
		MinionJob{
			Target:   ExpressionTarget{Expression: "minion1", Type: Glob},
			Function: "test.ping",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Contains(t, res[0].Minions, "minion1")
	assert.NotEmpty(t, res[0].ID)
	assert.Contains(t, res[1].Minions, "minion1")
	assert.NotEmpty(t, res[1].ID)
}

func TestSubmitMultipleJobToOfflineMinion(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_submit", "multiple_offline")

	res, err := c.SubmitJobs(context.Background(), []MinionJob{
		MinionJob{
			Target:   ExpressionTarget{Expression: "minion2", Type: Glob},
			Function: "test.ping",
		},
		MinionJob{
			Target:   ExpressionTarget{Expression: "minion2", Type: Glob},
			Function: "test.ping",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Contains(t, res[0].Minions, "minion2")
	assert.NotEmpty(t, res[0].ID)
	assert.Contains(t, res[1].Minions, "minion2")
	assert.NotEmpty(t, res[1].ID)
}

func TestSubmitMuiltipleJobsToMissingMinion(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_submit", "multiple_missing")

	res, err := c.SubmitJobs(context.Background(), []MinionJob{
		MinionJob{
			Target:   ExpressionTarget{Expression: "minion3", Type: Glob},
			Function: "test.ping",
		},
		MinionJob{
			Target:   ExpressionTarget{Expression: "minion3", Type: Glob},
			Function: "test.ping",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.Empty(t, res[0].ID)
	assert.Empty(t, res[1].ID)
}

func TestSubmitMuiltipleJobsToMixedStatusMinions(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "minions_submit", "multiple_mixed")

	res, err := c.SubmitJobs(context.Background(), []MinionJob{
		MinionJob{
			Target:   ExpressionTarget{Expression: "minion1", Type: Glob},
			Function: "test.ping",
		},
		MinionJob{
			Target:   ExpressionTarget{Expression: "minion2", Type: Glob},
			Function: "test.ping",
		},
		MinionJob{
			Target:   ExpressionTarget{Expression: "minion3", Type: Glob},
			Function: "test.ping",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, len(res))
	assert.Contains(t, res[0].Minions, "minion1")
	assert.NotEmpty(t, res[0].ID)
	assert.Contains(t, res[1].Minions, "minion2")
	assert.NotEmpty(t, res[1].ID)
	assert.Empty(t, res[2].Minions)
	assert.Empty(t, res[2].ID)
}
