package cherrypy

// Stats retrieves CherryPy stats
func (c *Client) Stats() (map[string]interface{}, error) {
	return c.requestJSON("GET", "stats", nil)
}
