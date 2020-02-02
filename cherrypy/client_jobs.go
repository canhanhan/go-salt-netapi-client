package cherrypy

import (
	"errors"
	"fmt"
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

/*
Job retrieves details of a single job

If the job was not found; ErrorJobNotFound will be returned.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--jobs-(jid)
*/
func (c *Client) Job(id string) (*JobDetails, error) {
	res, err := c.requestJSON("GET", "jobs/"+id, nil)
	if err != nil {
		return nil, err
	}

	infoBlock := res["info"].([]interface{})
	if len(infoBlock) != 1 {
		return nil, fmt.Errorf("expected 1 result in info block but received %d results", len(infoBlock))
	}

	dict := infoBlock[0].(map[string]interface{})
	if _, ok := dict["jid"]; !ok {
		return nil, fmt.Errorf("%s: %w", id, ErrorJobNotFound)
	}
	if v, ok := dict["Error"]; ok {
		return nil, errors.New(v.(string))
	}

	startTime, err := parseTime(dict["StartTime"].(string))
	if err != nil {
		return nil, err
	}

	job := JobDetails{}
	job.ID = dict["jid"].(string)
	job.Function = dict["Function"].(string)
	job.StartTime = startTime
	job.User = dict["User"].(string)
	job.Minions = stringSlice(dict["Minions"].([]interface{}))
	job.Returns = dict["Result"].(map[string]interface{})
	job.Arguments, job.KWArguments = parseArgs(dict["Arguments"].([]interface{}))

	targetType := targetTypes[dict["Target-type"].(string)]
	if targetType == List {
		job.Target = &ListTarget{
			Targets: stringSlice(dict["Target"].([]interface{})),
		}
	} else {
		job.Target = &ExpressionTarget{
			Expression: dict["Target"].(string),
			Type:       targetType,
		}
	}

	return &job, nil
}

/*
Jobs retrieves status of all jobs from Salt Master.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#get--jobs-(jid)
*/
func (c *Client) Jobs() ([]Job, error) {
	res, err := c.requestJSON("GET", "jobs", nil)
	if err != nil {
		return nil, err
	}

	results := res["return"].([]interface{})
	if len(results) != 1 {
		return nil, fmt.Errorf("expected 1 result but received %d results", len(results))
	}

	dict := results[0].(map[string]interface{})
	jobs := make([]Job, len(dict))
	i := 0
	for k, v := range dict {
		j := v.(map[string]interface{})
		startTime, err := parseTime(j["StartTime"].(string))
		if err != nil {
			return nil, err
		}

		job := Job{
			ID:        k,
			Function:  j["Function"].(string),
			StartTime: startTime,
			User:      j["User"].(string),
		}

		targetType := targetTypes[j["Target-type"].(string)]
		if targetType == List {
			job.Target = &ListTarget{
				Targets: stringSlice(j["Target"].([]interface{})),
			}
		} else {
			job.Target = &ExpressionTarget{
				Expression: j["Target"].(string),
				Type:       targetType,
			}
		}
		job.Arguments, job.KWArguments = parseArgs(j["Arguments"].([]interface{}))

		jobs[i] = job
		i++
	}

	return jobs, nil
}

func parseTime(val string) (time.Time, error) {
	return time.Parse("2006, Jan 02 15:04:05.000000", val)
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
