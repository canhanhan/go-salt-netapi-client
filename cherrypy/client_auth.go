package cherrypy

import (
	"context"
	"errors"
	"log"
)

var (
	// ErrorNotAuthenticated indicates Logout() was called before authenticating with Salt
	ErrorNotAuthenticated = errors.New("not authenticated")
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Backend  string `json:"eauth"`
}

type loginData struct {
	// TODO: Not sure about the data structure of permissions
	//	Permissions []string `json:"perms"`
	StartTime  saltUnixTime `json:"start"`
	Token      string       `json:"token"`
	ExpireTime saltUnixTime `json:"expire"`
	User       string       `json:"user"`
	Backend    string       `json:"eauth"`
}

type loginResponse struct {
	Return []loginData `json:"return"`
}

/*
Login establishes a session with rest_cherrypy and retrieves the token

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#login
*/
func (c *Client) Login(ctx context.Context) error {
	data := loginRequest{
		Username: c.eauth.Username,
		Password: c.eauth.Password,
		Backend:  c.eauth.Backend,
	}

	req, err := c.newRequest(ctx, "POST", "login", data)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] Sending authentication request")
	var response loginResponse
	_, err = c.do(req, &response)
	if err != nil {
		return err
	}

	c.Token = response.Return[0].Token
	log.Printf("[DEBUG] Received token %s", c.Token)

	return nil
}

/*
Logout terminates the session with rest_cherrypy and clears the token

Calls to logout will fail with ErrorNotAuthenticated if Login() was not called prior.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#logout
*/
func (c *Client) Logout(ctx context.Context) error {
	if c.Token == "" {
		return ErrorNotAuthenticated
	}

	req, err := c.newRequest(ctx, "POST", "logout", nil)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] Sending logout request")
	_, err = c.do(req, nil)
	if err != nil {
		return err
	}

	c.Token = ""
	return nil
}
