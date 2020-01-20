package cherrypy

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestValidLogin(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/login", "auth_login_success")

	err := c.Login()

	assert.NoError(t, err)
	assert.Equal(t, testToken, c.Token)
}

func TestInvalidLogin(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Some error", 401)
	})

	err := c.Login()

	assert.Error(t, err)
	assert.Equal(t, "", c.Token)
}

func TestLogout(t *testing.T) {
	c, mux, teardown := setup(t)
	defer teardown()
	handleJSONRequest(mux, "/login", "auth_login_success")
	handleJSONRequest(mux, "/logout", "auth_logout_success")
	if err := c.Login(); err != nil {
		t.Fatal(err)
	}

	err := c.Logout()

	assert.NoError(t, err)
	assert.Empty(t, c.Token)
}
