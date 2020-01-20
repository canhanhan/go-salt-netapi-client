package cherrypy

import (
	"fmt"
	"log"
)

// Login establishes a session with CherryPy and retrieves token
func (c *Client) Login() error {
	data := make(map[string]interface{})
	data["username"] = c.EAuth.Username
	data["password"] = c.EAuth.Password
	data["eauth"] = c.EAuth.Backend

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

// Logout terminates the session with CherryPy and clears the token
func (c *Client) Logout() error {
	if c.Token == "" {
		return fmt.Errorf("not authenticated")
	}

	log.Println("[DEBUG] Sending logout request")
	if _, err := c.requestJSON("POST", "logout", nil); err != nil {
		return err
	}

	c.Token = ""
	return nil
}
