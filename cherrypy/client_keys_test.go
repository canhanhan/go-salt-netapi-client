package cherrypy

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetKeys(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/keys", "key_list_success")

	res, err := c.Keys()

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "master.pem", res.Local[0])
	assert.Equal(t, "master.pub", res.Local[1])
	assert.Empty(t, res.MinionsRejected)
	assert.Empty(t, res.MinionsDenied)
	assert.Empty(t, res.MinionsPre)
	assert.Equal(t, "minion1", res.Minions[0])
	assert.Equal(t, "minion2", res.Minions[1])
}

func TestGetKey(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/keys/minion1", "key_get_success")

	res, err := c.Key("minion1")

	assert.NoError(t, err)
	assert.Equal(t, "b2:96:7c:28:2a:91:0a:7f:7a:8e:de:c1:dd:dd:cc:83:49:4f:ab:a9:a8:91:f8:80:19:2b:b8:e1:ec:9b:e5:57", res)
}

func TestGetKeyMissingMinion(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/keys/minion3", "key_get_missing")

	_, err := c.Key("minion3")

	assert.Error(t, err)
}

func TestGenerateKeySuccess(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()

	mux.HandleFunc("/keys", func(w http.ResponseWriter, req *http.Request) {
		content, err := getResponse("key_generate_success.bin")
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Write(content)
		}
	})

	res, err := c.GenerateKeyPair("minion1", 2048, false)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Public)
	assert.NotEmpty(t, res.Private)
}

func TestGenerateKeyFailure(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()

	mux.HandleFunc("/keys", func(w http.ResponseWriter, req *http.Request) {
		content, err := getResponse("key_generate_failure.bin")
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Write(content)
		}
	})

	_, err := c.GenerateKeyPair("minion1", 2048, false)

	assert.Error(t, err)
}
