package cherrypy

type KeyResult struct {
	Local           []string
	MinionsRejected []string
	MinionsDenied   []string
	MinionsPre      []string
	Minions         []string
}

// GetKeys does something
func (c *Client) GetKeys() (*KeyResult, error) {
	res, err := c.sendRequest("GET", "keys", nil)
	if err != nil {
		return nil, err
	}

	result := res["return"].(map[string]interface{})
	return &KeyResult{
		Local:           makeStringArray(result["local"].([]interface{})),
		MinionsRejected: makeStringArray(result["minions_rejected"].([]interface{})),
		MinionsDenied:   makeStringArray(result["minions_denied"].([]interface{})),
		MinionsPre:      makeStringArray(result["minions_pre"].([]interface{})),
		Minions:         makeStringArray(result["minions"].([]interface{})),
	}, nil
}

func makeStringArray(values []interface{}) []string {
	strings := make([]string, len(values))
	for i, v := range values {
		strings[i] = v.(string)
	}

	return strings
}
