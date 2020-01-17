package cherrypy

import (
	"fmt"
)

type Minion struct {
	ID     string
	Grains map[string]interface{}
}

// GetMinions does something
func (c *Client) GetMinions() ([]Minion, error) {
	res, err := c.sendRequest("GET", "minions", nil)
	if err != nil {
		return nil, err
	}

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
