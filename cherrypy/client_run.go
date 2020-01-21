package cherrypy

import (
	"fmt"
	"log"
)

/*
CommandClient indicates Salt API which client to use

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#usage

See the constants available in this file for possible values.
*/
type CommandClient string

const (
	// LocalClient sends commands to Minions. Equivalent to the salt CLI command.
	LocalClient CommandClient = "local"

	// RunnerClient invokes runner modules on the Master.
	// Equivalent to the salt-run CLI command.
	RunnerClient = "runner"

	// WheelClient invokes wheel modules on the Master.
	// Wheel modules do not have a direct CLI equivalent
	WheelClient = "wheel"
)

// Command to send to Run endpont
type Command struct {
	Client     CommandClient
	Target     string
	TargetType TargetType
	Function   string
	Args       []string
	Kwargs     map[string]interface{}
}

/*
RunJob runs a command on master using Run endpoint

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Run
*/
func (c *Client) RunJob(cmd Command) (map[string]interface{}, error) {
	res, err := c.RunJobs([]Command{cmd})
	if err != nil {
		return nil, err
	}

	if len(res) != 1 {
		return nil, fmt.Errorf("expected 1 results but received %d", len(res))
	}

	return res[0], nil
}

/*
RunJobs runs multiple commands on master using Run endpoint

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Run
*/
func (c *Client) RunJobs(cmds []Command) ([]map[string]interface{}, error) {
	items := make([]interface{}, len(cmds))
	for i, cmd := range cmds {
		data := make(map[string]interface{})
		data["client"] = cmd.Client
		if cmd.Target != "" {
			data["tgt"] = cmd.Target
		}
		if cmd.TargetType != "" {
			data["tgt_type"] = cmd.TargetType
		}
		if cmd.Function != "" {
			data["fun"] = cmd.Function
		}
		if len(cmd.Args) > 0 {
			data["arg"] = cmd.Args
		}
		if cmd.Kwargs != nil {
			data["kwarg"] = cmd.Kwargs
		}

		data["username"] = c.eauth.Username
		data["password"] = c.eauth.Password
		data["eauth"] = c.eauth.Backend

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
