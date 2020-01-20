package cherrypy

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type JobTarget interface {
	Target() string
	Type() string
}

type Job struct {
	ID          string
	Function    string
	Target      JobTarget
	Arguments   []interface{}
	KWArguments map[string]interface{}
	StartTime   time.Time
	User        string
}

type JobDetails struct {
	Job
	Minions []string
	Returns map[string]interface{}
}

type GlobTargetType struct {
	Expression string
}

func (c GlobTargetType) Target() string {
	return c.Expression
}
func (c GlobTargetType) Type() string {
	return "glob"
}

type ListTargetType struct {
	Minions []string
}

func (c ListTargetType) Target() string {
	return strings.Join(c.Minions, ", ")
}
func (c ListTargetType) Type() string {
	return "list"
}

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
		return nil, fmt.Errorf("job %s was not found", id)
	}
	if v, ok := dict["Error"]; ok {
		return nil, errors.New(v.(string))
	}

	startTime, err := parseTime(dict["StartTime"].(string))
	if err != nil {
		return nil, err
	}

	targetType := dict["Target-type"].(string)
	target, err := parseTarget(dict["Target"].(string), targetType)
	if err != nil {
		return nil, err
	}

	job := JobDetails{}
	job.ID = dict["jid"].(string)
	job.Function = dict["Function"].(string)
	job.StartTime = startTime
	job.Target = target
	job.User = dict["User"].(string)
	job.Minions = stringSlice(dict["Minions"].([]interface{}))
	job.Returns = dict["Result"].(map[string]interface{})

	job.Arguments, job.KWArguments = parseArgs(dict["Arguments"].([]interface{}))

	return &job, nil
}

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

		targetType := j["Target-type"].(string)
		target, err := parseTarget(j["Target"], targetType)
		if err != nil {
			return nil, err
		}

		job := Job{
			ID:        k,
			Function:  j["Function"].(string),
			StartTime: startTime,
			Target:    target,
			User:      j["User"].(string),
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

func parseTarget(target interface{}, targetType string) (JobTarget, error) {
	switch targetType {
	case "glob":
		return GlobTargetType{
			Expression: target.(string),
		}, nil
	case "list":
		return ListTargetType{
			Minions: stringSlice(target.([]interface{})),
		}, nil
	default:
		return nil, fmt.Errorf("unknown target type: %s", targetType)
	}
}
