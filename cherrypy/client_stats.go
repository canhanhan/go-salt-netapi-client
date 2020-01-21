package cherrypy

/*
Stats retrieves CherryPy stats

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#stats
*/
func (c *Client) Stats() (map[string]interface{}, error) {
	return c.requestJSON("GET", "stats", nil)
}
