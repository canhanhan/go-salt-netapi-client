package cherrypy

func stringSlice(raw []interface{}) []string {
	x := make([]string, len(raw))
	for i, v := range raw {
		x[i] = v.(string)
	}

	return x
}
