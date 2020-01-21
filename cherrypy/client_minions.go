package cherrypy

import (
	"errors"
	"fmt"
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
	Target     string
	TargetType TargetType
	Function   string
	Args       []interface{}
	KWArgs     map[string]interface{}
}

// AsyncMinionJobResult contains results of an async run with local client.
type AsyncMinionJobResult struct {
	Minions []string
	JobID   string
}

/*
Minion retrieves grains of a single minion from Salt Master

If the minion is offline; grains will be empty.
If the requested minion is not known by the master; ErrorMinionNotFound error will be thrown.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--minions-(mid)
*/
func (c *Client) Minion(id string) (*Minion, error) {
	minions, err := c.getMinions(id)

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
func (c *Client) Minions() ([]Minion, error) {
	return c.getMinions("")
}

/*
SubmitJobs submits multiple jobs to be executed on minions asynchronously

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Minions.POST
*/
func (c *Client) SubmitJobs(jobs []MinionJob) ([]AsyncMinionJobResult, error) {
	data := make([]interface{}, len(jobs))
	for i, v := range jobs {
		job := make(map[string]interface{})
		job["tgt"] = v.Target
		job["tgt_type"] = v.TargetType
		job["fun"] = v.Function
		job["args"] = v.Args
		job["kwargs"] = v.KWArgs

		data[i] = job
	}

	res, err := c.requestJSON("POST", "minions", data)
	if err != nil {
		return nil, err
	}

	rawResults := res["return"].([]interface{})
	results := make([]AsyncMinionJobResult, 0)
	for _, v := range rawResults {
		dict := v.(map[string]interface{})
		// If target did not match to any minions Salt returns empty object. Skip...
		if len(dict) == 0 {
			continue
		}

		results = append(results, AsyncMinionJobResult{
			Minions: stringSlice(dict["minions"].([]interface{})),
			JobID:   dict["jid"].(string),
		})
	}

	return results, nil
}

/*
SubmitJob submits a single job to be executed on minions asynchronously

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Minions.POST
*/
func (c *Client) SubmitJob(job MinionJob) (*AsyncMinionJobResult, error) {
	res, err := c.SubmitJobs([]MinionJob{job})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

func (c *Client) getMinions(id string) ([]Minion, error) {
	res, err := c.requestJSON("GET", "minions/"+id, nil)
	if err != nil {
		return nil, err
	}

	r := res["return"].([]interface{})
	if len(r) != 1 {
		return nil, fmt.Errorf("expected one return but received %d", len(r))
	}

	dict := r[0].(map[string]interface{})
	minions := make([]Minion, len(dict))

	i := 0
	for k, m := range dict {
		minions[i] = Minion{ID: k}

		// Grains are not returned for offline minions
		if g, ok := m.(map[string]interface{}); ok {
			minions[i].Grains = g
		}

		i++
	}

	return minions, nil
}
