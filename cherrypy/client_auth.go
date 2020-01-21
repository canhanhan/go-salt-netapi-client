package cherrypy

import (
	"errors"
	"log"
)

var (
	// ErrorNotAuthenticated indicates Logout() was called before authenticating with Salt
	ErrorNotAuthenticated = errors.New("not authenticated")
)

/*
Login establishes a session with rest_cherrypy and retrieves the token

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#login
*/
func (c *Client) Login() error {
	data := make(map[string]interface{})
	data["username"] = c.eauth.Username
	data["password"] = c.eauth.Password
	data["eauth"] = c.eauth.Backend

	log.Println("[DEBUG] Sending authentication request")
	result, err := c.requestJSON("POST", "login", data)
	if err != nil {
		return err
	}

	records := result["return"].([]interface{})
	record := records[0].(map[string]interface{})
	c.Token = record["token"].(string)

	log.Printf("[DEBUG] Received token %s", c.Token)

	return nil
}

/*
Logout terminates the session with rest_cherrypy and clears the token

Calls to logout will fail with ErrorNotAuthenticated if Login() was not called prior.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#logout
*/
func (c *Client) Logout() error {
	if c.Token == "" {
		return ErrorNotAuthenticated
	}

	log.Println("[DEBUG] Sending logout request")
	if _, err := c.requestJSON("POST", "logout", nil); err != nil {
		return err
	}

	c.Token = ""
	return nil
}
