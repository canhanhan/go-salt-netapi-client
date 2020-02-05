package cherrypy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeys(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "keys_list", "success")

	res, err := c.Keys(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "master.pem", res.Local[0])
	assert.Equal(t, "master.pub", res.Local[1])
	assert.Empty(t, res.MinionsRejected)
	assert.Empty(t, res.MinionsDenied)
	assert.Equal(t, "saltmaster.local", res.MinionsPre[0])
	assert.Equal(t, "minion1", res.Minions[0])
	assert.Equal(t, "minion2", res.Minions[1])
}

func TestGetKey(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "keys_get", "success")

	res, err := c.Key(context.Background(), "minion1")

	assert.NoError(t, err)
	assert.Equal(t, "b2:96:7c:28:2a:91:0a:7f:7a:8e:de:c1:dd:dd:cc:83:49:4f:ab:a9:a8:91:f8:80:19:2b:b8:e1:ec:9b:e5:57", res)
}

func TestGetKeyMissingMinion(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "keys_get", "missing")

	_, err := c.Key(context.Background(), "minion3")

	assert.Error(t, err)
}

func TestGenerateKeySuccess(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "keys_generate", "success")

	res, err := c.GenerateKeyPair(context.Background(), "minion4", 2048, false)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Public)
	assert.NotEmpty(t, res.Private)
}

func TestGenerateKeyFailure(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "keys_generate", "failure")

	_, err := c.GenerateKeyPair(context.Background(), "minion4", 2048, false)

	assert.Error(t, err)
}
