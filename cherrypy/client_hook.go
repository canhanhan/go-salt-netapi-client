package cherrypy

import (
	"context"
	"fmt"
	"log"
)

type hookResponse struct {
	Success bool   `json:"success,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

/*
Hook fires an event on Salt's event bus

All events are prefixed with salt/netapi/hook.
Therefore if the id is set to "test"; Salt Reactor will receive "salt/netapi/hook/test" event.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Webhook.POST
*/
func (c *Client) Hook(ctx context.Context, id string, data interface{}) error {
	req, err := c.newRequest(ctx, "POST", "hook/"+id, data)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] Sending authentication request")
	var resp hookResponse
	_, err = c.do(req, &resp)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("unexpected response from Salt: %d, %s", resp.Status, resp.Message)
	}

	return nil
}
