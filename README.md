# go-salt-netapi-client
[![GoDoc](https://godoc.org/github.com/finarfin/go-salt-netapi-client/cherrypy?status.svg)](https://godoc.org/finarfin/go-salt-netapi-client/cherrypy)
[![Test Status](https://github.com/finarfin/go-salt-netapi-client/workflows/Go/badge.svg)](https://github.com/finarfin/go-salt-netapi-client/actions?query=workflow%3AGo)

go-salt-netapi-client is a Go client library for accessing the [NetAPI modules](https://docs.saltstack.com/en/latest/ref/netapi/all/index.html) of [SaltStack OSS](https://github.com/saltstack/salt). Currently only [rest_cherrypy](https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html) is supported.

go-salt-netapi-client requires Go version 1.13 or greater.

## Usage ##

```go
import "github.com/finarfin/go-salt-netapi-client/cherrypy"
```

Construct a new client, then use the various methods on the client.

```go
client := cherrypy.NewClient("https://master:8000", "admin", "password", "pam")

// list all minions
minions, err := client.Minions()
```

See [GoDoc](https://godoc.org/finarfin/go-salt-netapi-client/cherrypy) for details.
