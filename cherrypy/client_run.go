package cherrypy

import (
	"fmt"
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

// RunJob runs a command on master using Run endpoint
func (c *Client) RunJob(cmd Command) (map[string]interface{}, error) {
	res, err := c.RunJobs([]Command { cmd })
	if err != nil {
		return nil, err
	}

	if len(res) != 1 {
		return nil, fmt.Errorf("expected 1 results but received %d", len(res))
	}

	return res[0], nil
}

// RunJobs runs multiple commands on master using Run endpoint
func (c *Client) RunJobs(cmds []Command) ([]map[string]interface{}, error) {
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
	res, err := c.requestJSON("POST", "run", items)
	if err != nil {
		return nil, err
	}
	
	results := res["return"].([]interface{})
	output := make([]map[string]interface{}, len(results))
	for i, v := range results {
		output[i] = v.(map[string]interface{})
	}

	return output, nil
}

