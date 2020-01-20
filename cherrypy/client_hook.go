package cherrypy

import (
	"fmt"
)

// Hook fires an event on Salt's event bus
func (c *Client) Hook(id string, data map[string]interface{}) error {
	res, err := c.requestJSON("POST", "hook/"+id, data)
	if err != nil {
		return err
	}

	// CherryPy returns {{ success: true }} if the hook is received
	if successRaw, ok := res["success"]; ok {
		if success, ok := successRaw.(bool); !ok || !success {
			return fmt.Errorf("unexpected status: %v", successRaw)
		}

		return nil
	}

	return fmt.Errorf("unexpected response from Salt: %v", res)
}
