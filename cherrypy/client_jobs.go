package cherrypy

import (
	"context"
	"errors"
	"log"
	"time"
)

var (
	// ErrorJobNotFound indicates requested job was not found
	ErrorJobNotFound = errors.New("job was not found")
)

// Job contains summary of a job returned by Jobs()
type Job struct {
	ID          string
	Function    string
	Target      Target
	Arguments   []interface{}
	KWArguments map[string]interface{}
	StartTime   time.Time
	User        string
}

// JobDetails contain job summary and returns per minion
type JobDetails struct {
	Job
	Minions []string
	Returns map[string]interface{}
}

type jobResult struct {
	Return     interface{} `json:"return"`
	ReturnCode int         `json:"retcode"`
	Success    bool        `json:"success"`
}

type jobInfo struct {
	Function   string               `json:"Function"`
	ID         string               `json:"jid,omitempty"`
	Result     map[string]jobResult `json:"Result"`
	User       string               `json:"User"`
	Target     interface{}          `json:"Target"`
	TargetType string               `json:"Target-type"`
	StartTime  saltTime             `json:"StartTime"`
	Minions    []string             `json:"Minions"`
	Arguments  []interface{}        `json:"Arguments"`
}

type jobDetailResponse struct {
	Info    []jobInfo                `json:"info"`
	Returns []map[string]interface{} `json:"return"`
}

type jobListResponse struct {
	Jobs []map[string]jobInfo `json:"return"`
}

/*
Job retrieves details of a single job

If the job was not found; ErrorJobNotFound will be returned.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--jobs-(jid)
*/
func (c *Client) Job(ctx context.Context, id string) (*JobDetails, error) {
	req, err := c.newRequest(ctx, "GET", "jobs/"+id, nil)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Sending job details request")
	var resp jobDetailResponse
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}

	j := resp.Info[0]
	job := JobDetails{
		Minions: j.Minions,
		Returns: resp.Returns[0],
	}

	job.ID = j.ID
	job.Function = j.Function
	job.StartTime = j.StartTime.Time
	job.User = j.User
	job.Arguments, job.KWArguments = parseArgs(j.Arguments)
	job.Target = parseTarget(j)

	return &job, nil
}

/*
Jobs retrieves status of all jobs from Salt Master.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--jobs-(jid)
*/
func (c *Client) Jobs(ctx context.Context) ([]Job, error) {
	req, err := c.newRequest(ctx, "GET", "jobs", nil)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Sending job list request")
	var resp jobListResponse
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}

	jobs := make([]Job, len(resp.Jobs[0]))
	i := 0
	for k, v := range resp.Jobs[0] {
		args, kwArgs := parseArgs(v.Arguments)
		target := parseTarget(v)

		jobs[i] = Job{
			ID:          k,
			Function:    v.Function,
			StartTime:   v.StartTime.Time,
			User:        v.User,
			Target:      target,
			Arguments:   args,
			KWArguments: kwArgs,
		}

		i++
	}

	return jobs, nil
}

func parseTarget(j jobInfo) Target {
	targetType := targetTypes[j.TargetType]
	switch targetType {
	case List:
		return &ListTarget{
			Targets: stringSlice(j.Target.([]interface{})),
		}
	default:
		return &ExpressionTarget{
			Expression: j.Target.(string),
			Type:       targetType,
		}
	}
}

func parseArgs(arguments []interface{}) ([]interface{}, map[string]interface{}) {
	args := make([]interface{}, 0)
	kwargs := make(map[string]interface{})
	for _, arg := range arguments {
		if d, ok := arg.(map[string]interface{}); ok {
			if kw, ok := d["__kwarg__"]; ok && kw.(bool) {
				for k, v := range d {
					if k == "__kwarg__" {
						continue
					}
					kwargs[k] = v
				}
			} else {
				args = append(args, d)
			}
		} else {
			args = append(args, arg)
		}
	}

	return args, kwargs
}
