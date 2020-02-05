package cherrypy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHookSuccess(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "hook", "success")

	err := c.Hook(context.Background(), "test", nil)

	assert.NoError(t, err)
}

func TestHookFailure(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "hook", "failure")

	err := c.Hook(context.Background(), "test", nil)

	assert.Error(t, err)
}
