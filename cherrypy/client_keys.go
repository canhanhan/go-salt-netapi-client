package cherrypy

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	Local           []string `json:"local"`
	MinionsRejected []string `json:"minions_rejected"`
	MinionsDenied   []string `json:"minions_denied"`
	MinionsPre      []string `json:"minions_pre"`
	Minions         []string `json:"minions"`
}

// MinionKeyPair contains key-pair for a minion
type MinionKeyPair struct {
	ID      string
	Public  string
	Private string
}

type keyListResponse struct {
	Return KeyResult `json:"return"`
}

type keyDetailsMinionsKeyPair struct {
	Minions map[string]string `json:"minions"`
}

type keyDetailsResponse struct {
	Return keyDetailsMinionsKeyPair `json:"return"`
}

type keyGenerateRequest struct {
	ID       string `json:"mid"`
	KeySize  int    `json:"keysize,omitempty"`
	Force    bool   `json:"force"`
	Username string `json:"username"`
	Password string `json:"password"`
	Backend  string `json:"eauth"`
}

/*
Keys retrieves list of keys from master

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Keys.GET
*/
func (c *Client) Keys(ctx context.Context) (*KeyResult, error) {
	req, err := c.newRequest(ctx, "GET", "keys", nil)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Sending key list request")
	var resp keyListResponse
	_, err = c.do(req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Return, nil
}

/*
Key returns public key of a single minion from master

If the minion is not found on the master ErrorMinionKeyNotFound error is returned.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Keys.GET
*/
func (c *Client) Key(ctx context.Context, id string) (string, error) {
	req, err := c.newRequest(ctx, "GET", "keys/"+id, nil)
	if err != nil {
		return "", err
	}

	log.Println("[DEBUG] Sending key details request")
	var resp keyDetailsResponse
	_, err = c.do(req, &resp)
	if err != nil {
		return "", err
	}

	if len(resp.Return.Minions) == 0 {
		return "", fmt.Errorf("%s: %w", id, ErrorMinionKeyNotFound)
	}

	return resp.Return.Minions[id], nil
}

/*
GenerateKeyPair generates and auto-accepts minion keypair on the master.

If force argument is false and if the master already has keys for the minion;
ErrorKeyPairNotReceived error will be returned.

If force argument is true; existing keys will be overwriten and new keys will be generated.

https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html#salt.netapi.rest_cherrypy.app.Keys.POST
*/
func (c *Client) GenerateKeyPair(ctx context.Context, id string, keySize int, force bool) (*MinionKeyPair, error) {
	data := keyGenerateRequest{
		ID:       id,
		KeySize:  keySize,
		Force:    force,
		Username: c.eauth.Username,
		Password: c.eauth.Password,
		Backend:  c.eauth.Backend,
	}

	req, err := c.newRequest(ctx, "POST", "keys", data)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Sending generate key request")
	br := new(bytes.Buffer)
	_, err = c.do(req, br)
	if err != nil {
		return nil, err
	}

	keys := MinionKeyPair{
		ID: id,
	}

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
