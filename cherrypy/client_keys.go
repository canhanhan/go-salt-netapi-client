package cherrypy

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

type KeyResult struct {
	Local           []string
	MinionsRejected []string
	MinionsDenied   []string
	MinionsPre      []string
	Minions         []string
}

type MinionKeys struct {
	Public  string
	Private string
}

// Keys retrieves list of keys from master
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

// Key returns minion public key from master
func (c *Client) Key(id string) (string, error) {
	res, err := c.requestJSON("GET", "keys/"+id, nil)
	if err != nil {
		return "", err
	}

	result := res["return"].(map[string]interface{})
	if len(result) == 0 {
		return "", errors.New("key was not found")
	}

	dict := result["minions"].(map[string]interface{})
	return dict[id].(string), nil
}

func (c *Client) GenerateKeys(id string, keySize int, force bool) (*MinionKeys, error) {
	data := make(map[string]interface{})
	data["mid"] = id
	data["keysize"] = strconv.Itoa(keySize)
	data["force"] = strconv.FormatBool(force)
	data["username"] = c.EAuth.Username
	data["password"] = c.EAuth.Password
	data["eauth"] = c.EAuth.Backend

	body, err := c.request("POST", "keys", data)
	if err != nil {
		return nil, err
	}

	keys := MinionKeys{}

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
			return nil, fmt.Errorf("unexpected file type in tar file for: %s", target)
		}

		content, err := ioutil.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("error reading contents of %s", target)
		}

		switch target {
		case "minion.pub":
			keys.Public = string(content)
		case "minion.pem":
			keys.Private = string(content)
		default:
			return nil, fmt.Errorf("unknown file name in tar file: %s", target)
		}
	}

	if keys.Public == "" || keys.Private == "" {
		return nil, errors.New("public or private key was not received")
	}

	return &keys, nil
}
