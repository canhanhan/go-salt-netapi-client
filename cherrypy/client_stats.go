package cherrypy

import (
	"context"
	"log"
)

/*
Stats retrieves CherryPy stats

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#stats
*/
func (c *Client) Stats(ctx context.Context) (map[string]interface{}, error) {
	req, err := c.newRequest(ctx, "GET", "stats", nil)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Sending stats request")
	var resp map[string]interface{}
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
