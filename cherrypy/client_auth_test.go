package cherrypy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidLogin(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "auth_login", "success")

	c.Token = ""
	err := c.Login(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, testToken, c.Token)
}

func TestInvalidLogin(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "auth_login", "bad_user")

	c.Token = ""
	err := c.Login(context.Background())

	assert.Error(t, err)
	assert.Equal(t, "", c.Token)
}

func TestLogout(t *testing.T) {
	tester, c := setup(t)
	defer tester.Close()
	tester.Setup(t, "auth_login", "success")
	tester.Setup(t, "auth_logout", "success")

	err := c.Login(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	err = c.Logout(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, c.Token)
}
