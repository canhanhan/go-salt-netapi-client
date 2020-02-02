package cherrypy

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

var (
	// ErrorMinionKeyNotFound indicates requested minion key does not exist on the master
	ErrorMinionKeyNotFound = errors.New("minion key was not found")

	// ErrorKeyPairNotReceived indicates master did not create a new key-pair.
	// This occurs if a key was already accepted by master for that particular minion
	// and the force argument was set to false
	ErrorKeyPairNotReceived = errors.New("public or private key was not received")

	// ErrorKeyPairBroken indicates the key-pair file received was broken
	ErrorKeyPairBroken = errors.New("cannot extract keypair from archive")
)

// KeyResult contains list of available keys on the master
type KeyResult struct {
	Local           []string
	MinionsRejected []string
	MinionsDenied   []string
	MinionsPre      []string
	Minions         []string
}

// MinionKeyPair contains key-pair for a minion
type MinionKeyPair struct {
	ID      string
	Public  string
	Private string
}

/*
Keys retrieves list of keys from master

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Keys.GET
*/
func (c *Client) Keys() (*KeyResult, error) {
	res, err := c.requestJSON("GET", "keys", nil)
	if err != nil {
		return nil, err
	}

	result := res["return"].(map[string]interface{})
	return &KeyResult{
		Local:           stringSlice(result["local"].([]interface{})),
		MinionsRejected: stringSlice(result["minions_rejected"].([]interface{})),
		MinionsDenied:   stringSlice(result["minions_denied"].([]interface{})),
		MinionsPre:      stringSlice(result["minions_pre"].([]interface{})),
		Minions:         stringSlice(result["minions"].([]interface{})),
	}, nil
}

/*
Key returns public key of a single minion from master

If the minion is not found on the master ErrorMinionKeyNotFound error is returned.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Keys.GET
*/
func (c *Client) Key(id string) (string, error) {
	res, err := c.requestJSON("GET", "keys/"+id, nil)
	if err != nil {
		return "", err
	}

	result := res["return"].(map[string]interface{})
	if len(result) == 0 {
		return "", fmt.Errorf("%s: %w", id, ErrorMinionKeyNotFound)
	}

	dict := result["minions"].(map[string]interface{})
	return dict[id].(string), nil
}

/*
GenerateKeyPair generates and auto-accepts minion keypair on the master.

If force argument is false and if the master already has keys for the minion;
ErrorKeyPairNotReceived error will be returned.

If force argument is true; existing keys will be overwriten and new keys will be generated.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Keys.POST
*/
func (c *Client) GenerateKeyPair(id string, keySize int, force bool) (*MinionKeyPair, error) {
	data := make(map[string]interface{})
	data["mid"] = id
	data["keysize"] = keySize
	data["force"] = force
	data["username"] = c.eauth.Username
	data["password"] = c.eauth.Password
	data["eauth"] = c.eauth.Backend

	body, err := c.request("POST", "keys", "", data)
	if err != nil {
		return nil, err
	}

	keys := MinionKeyPair{
		ID: id,
	}

	br := bytes.NewReader(body)
	tr := tar.NewReader(br)
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		target := header.Name
		if header.Typeflag != tar.TypeReg {
			return nil, fmt.Errorf("%s: %w unexpected file type", target, ErrorKeyPairBroken)
		}

		content, err := ioutil.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w error reading contents", target, ErrorKeyPairBroken)
		}

		switch target {
		case "minion.pub":
			keys.Public = string(content)
		case "minion.pem":
			keys.Private = string(content)
		default:
			return nil, fmt.Errorf("%s: %w unknown file name", target, ErrorKeyPairBroken)
		}
	}

	if keys.Public == "" || keys.Private == "" {
		return nil, fmt.Errorf("%s: %w", id, ErrorKeyPairNotReceived)
	}

	return &keys, nil
}
