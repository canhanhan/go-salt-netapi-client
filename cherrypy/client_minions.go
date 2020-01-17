package cherrypy

import (
	"fmt"
)

type Minion struct {
	ID     string
	Grains map[string]interface{}
}

func parseSomething(res map[string]interface{}) ([]Minion, error) {
	returns := res["return"].([]interface{})
	if len(returns) != 1 {
		return nil, fmt.Errorf("Expected 1 result but received %d", len(returns))
	}

	minionDict := returns[0].(map[string]interface{})
	minions := make([]Minion, len(minionDict))
	i := 0

	for k, v := range minionDict {
		minions[i] = Minion{
			ID: k,
		}

		if g, ok := v.(map[string]interface{}); ok {
			minions[i].Grains = g
		}

		i++
	}

	return minions, nil
}

// GetMinion does something
func (c *Client) GetMinion(minionID string) (*Minion, error) {
	res, err := c.sendRequest("GET", "minions/"+minionID, nil)
	if err != nil {
		return nil, err
	}

	minions, err := parseSomething(res)
	if err != nil {
		return nil, err
	}

	if len(minions) == 0 {
		return nil, nil
	}

	return &minions[0], nil
}

// GetMinions does something
func (c *Client) GetMinions() ([]Minion, error) {
	res, err := c.sendRequest("GET", "minions", nil)
	if err != nil {
		return nil, err
	}

	return parseSomething(res)
}
