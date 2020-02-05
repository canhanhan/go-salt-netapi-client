package cherrypy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunLocalCommand(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "run", "local_success")

	cmd := Command{
		Client:   "local",
		Target:   ExpressionTarget{Expression: "minion1", Type: Glob},
		Function: "test.ping",
	}

	res, err := c.RunJob(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestRunLocalCommands(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "run", "local_multiple_success")

	cmds := []Command{
		Command{
			Client:   "local",
			Target:   ExpressionTarget{Expression: "minion1", Type: Glob},
			Function: "test.ping",
		},
		Command{
			Client:   "local",
			Target:   ExpressionTarget{Expression: "minion1", Type: Glob},
			Function: "test.ping",
		},
	}

	res, err := c.RunJobs(context.Background(), cmds)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}

func TestRunWheelCommand(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "run", "wheel_success")

	cmd := Command{
		Client:   "wheel",
		Function: "minions.connected",
	}

	res, err := c.RunJob(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

// TODO: Add runner test
// TODO: Add test with arguments
// TODO: Add test with kw arguments
// TODO: Add tests with 401
