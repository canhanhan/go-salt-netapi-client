package cherrypy

import (
	"log"
)

// Command Command to send to Run endpont
type Command struct {
	Client   string
	Target   string
	Function string
	Args     []string
	Kwargs   map[string]string
}

// Run runs a command on master using Run endpoint
func (c *Client) Run(cmds []Command) (map[string]interface{}, error) {
	items := make([]interface{}, len(cmds))
	for i, cmd := range cmds {
		data := make(map[string]interface{})
		data["client"] = cmd.Client
		data["tgt"] = cmd.Target
		data["fun"] = cmd.Function
		data["arg"] = cmd.Args
		data["kwarg"] = cmd.Kwargs
		data["username"] = c.EAuth.Username
		data["password"] = c.EAuth.Password
		data["eauth"] = c.EAuth.Backend

		items[i] = data
	}

	log.Println("[DEBUG] Sending run request")
	return c.sendRequest("POST", "run", items)
}
