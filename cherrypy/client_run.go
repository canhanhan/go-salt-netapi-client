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
	Args       []string
	Kwargs     map[string]interface{}
	FullReturn bool
}

type runRequest struct {
	Client     CommandClient          `json:"client,omitempty"`
	Target     interface{}            `json:"tgt,omitempty"`
	TargetType TargetType             `json:"tgt_type,omitempty"`
	Function   string                 `json:"fun,omitempty"`
	Username   string                 `json:"username,omitempty"`
	Password   string                 `json:"password,omitempty"`
	Backend    string                 `json:"eauth,omitempty"`
	Args       []string               `json:"args,omitempty"`
	KWArgs     map[string]interface{} `json:"kwarg,omitempty"`
	FullReturn bool                   `json:"full_return,omitempty"`
}

type runResponse struct {
	Return []interface{} `json:"return"`
}

/*
RunJob runs a command on master using Run endpoint

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Run
*/
func (c *Client) RunJob(ctx context.Context, cmd Command) (interface{}, error) {
	res, err := c.RunJobs(ctx, []Command{cmd})
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
func (c *Client) RunJobs(ctx context.Context, cmds []Command) ([]interface{}, error) {
	d := make([]runRequest, len(cmds))
	for i, v := range cmds {
		j := runRequest{
			Args:     v.Args,
			KWArgs:   v.Kwargs,
			Client:   v.Client,
			Function: v.Function,
			Username: c.eauth.Username,
			Password: c.eauth.Password,
			Backend:  c.eauth.Backend,
		}

		if v.Target != nil {
			j.Target = v.Target.GetTarget()
			j.TargetType = v.Target.GetType()
		}

		// wheel throws following error if full_return is sent as a seperate argument
		// TypeError: call_func() got multiple values for keyword argument 'full_return'
		if j.Client == WheelClient && j.FullReturn {
			if j.KWArgs == nil {
				j.KWArgs = make(map[string]interface{})
			}
			j.KWArgs["full_return"] = j.FullReturn
			j.FullReturn = false
		}

		d[i] = j
	}

	req, err := c.newRequest(ctx, "POST", "run", d)
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
