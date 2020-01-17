package cherrypy

import (
	"log"
)

// Authenticate establishes a session with CherryPy
func (c *Client) Authenticate() error {
	data := make(map[string]interface{})
	data["username"] = c.EAuth.Username
	data["password"] = c.EAuth.Password
	data["eauth"] = c.EAuth.Backend

	log.Println("[DEBUG] Sending authentication request")
	result, err := c.sendRequest("POST", "login", data)
	if err != nil {
		return err
	}

	records := result["return"].([]interface{})
	record := records[0].(map[string]interface{})
	c.Token = record["token"].(string)

	log.Printf("[DEBUG] Received token %s", c.Token)

	return nil
}
