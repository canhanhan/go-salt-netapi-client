package cherrypy

import (
	"context"
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
	Target     Target
	Function   string
	Arguments  map[string]interface{}
}

type runResponse struct {
	Return []interface{} `json:"return"`
}

/*
RunCommand runs a command on master using Run endpoint

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Run
*/
func (c *Client) RunCommand(ctx context.Context, cmd Command) (interface{}, error) {
	res, err := c.RunCommands(ctx, []Command{cmd})
	if err != nil {
		return nil, err
	}

	if len(res) != 1 {
		return nil, fmt.Errorf("expected 1 results but received %d", len(res))
	}

	return res[0], nil
}

/*
RunCommands runs multiple commands on master using Run endpoint

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Run
*/
func (c *Client) RunCommands(ctx context.Context, cmds []Command) ([]interface{}, error) {
	r := make([]map[string]interface{}, len(cmds))
	for i, v := range cmds {
		d := make(map[string]interface{})
		
		if v.Arguments != nil {
			for k, a := range v.Arguments {
				d[k] = a
			}
		}

		d["client"] = v.Client
		d["fun"] = v.Function
		d["username"] = c.eauth.Username
		d["password"] = c.eauth.Password
		d["eauth"] = c.eauth.Backend

		if v.Target != nil {
			d["tgt"] = v.Target.GetTarget()
			d["tgt_type"] = v.Target.GetType()
		}

		// wheel throws following error if full_return is sent as a seperate argument
		// TypeError: call_func() got multiple values for keyword argument 'full_return'
		if v.Client != WheelClient {
			d["full_return"] = true
		}
		
		r[i] = d
	}

	req, err := c.newRequest(ctx, "POST", "run", r)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Sending run jobs request")
	var resp runResponse
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Return, nil
}
