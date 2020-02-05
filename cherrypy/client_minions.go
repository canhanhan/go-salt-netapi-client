package cherrypy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

var (
	// ErrorMinionNotFound indicates that minion was not found on Salt Master
	ErrorMinionNotFound = errors.New("minion not found")
)

// Minion information
type Minion struct {
	ID     string
	Grains map[string]interface{}
}

// MinionJob contains job information to be sent to the minion
type MinionJob struct {
	Target      Target
	Function    string
	Arguments   []interface{}
	KWArguments map[string]interface{}
}

// AsyncMinionJobResult contains results of an async run with local client.
type AsyncMinionJobResult struct {
	ID      string   `json:"jid"`
	Minions []string `json:"minions"`
}

type submitMinionJob struct {
	Target      interface{}            `json:"tgt"`
	TargetType  TargetType             `json:"tgt_type"`
	Function    string                 `json:"fun"`
	Arguments   []interface{}          `json:"args,omitempty"`
	KWArguments map[string]interface{} `json:"kwargs,omitempty"`
}

type submitMinionJobResponse struct {
	Return []AsyncMinionJobResult `json:"return"`
}

type minionDetailResponse struct {
	Return []map[string]json.RawMessage `json:"return"`
}

/*
Minion retrieves grains of a single minion from Salt Master

If the minion is offline; grains will be empty.
If the requested minion is not known by the master; ErrorMinionNotFound error will be thrown.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--minions-(mid)
*/
func (c *Client) Minion(ctx context.Context, id string) (*Minion, error) {
	minions, err := c.getMinions(ctx, id)

	if err != nil {
		return nil, err
	}

	if len(minions) == 0 {
		return nil, fmt.Errorf("%s: %w", id, ErrorMinionNotFound)
	}

	return &minions[0], nil
}

/*
Minions retrieves grains of all minions on a Salt Master

Grains will be empty for offline minions.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--minions-(mid)
*/
func (c *Client) Minions(ctx context.Context) ([]Minion, error) {
	return c.getMinions(ctx, "")
}

/*
SubmitJobs submits multiple jobs to be executed on minions asynchronously

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Minions.POST
*/
func (c *Client) SubmitJobs(ctx context.Context, jobs []MinionJob) ([]AsyncMinionJobResult, error) {
	data := make([]submitMinionJob, len(jobs))
	for i, v := range jobs {
		data[i] = submitMinionJob{
			Target:      v.Target.GetTarget(),
			TargetType:  v.Target.GetType(),
			Function:    v.Function,
			Arguments:   v.Arguments,
			KWArguments: v.KWArguments,
		}
	}

	req, err := c.newRequest(ctx, "POST", "minions", data)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Sending submit minion job request")
	var resp submitMinionJobResponse
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Return, nil
}

/*
SubmitJob submits a single job to be executed on minions asynchronously

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Minions.POST
*/
func (c *Client) SubmitJob(ctx context.Context, job MinionJob) (*AsyncMinionJobResult, error) {
	res, err := c.SubmitJobs(ctx, []MinionJob{job})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

func (c *Client) getMinions(ctx context.Context, id string) ([]Minion, error) {
	req, err := c.newRequest(ctx, "GET", "minions/"+id, nil)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Sending minion details request")
	var resp minionDetailResponse
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Return) != 1 {
		return nil, fmt.Errorf("expected one return but received %d", len(resp.Return))
	}

	d := resp.Return[0]
	minions := make([]Minion, len(d))

	i := 0
	for k, m := range d {
		minions[i] = Minion{ID: k}

		// Grains are not returned for offline minions
		var g map[string]interface{}
		if json.Unmarshal(m, &g) == nil {
			minions[i].Grains = g
		}

		i++
	}

	return minions, nil
}
