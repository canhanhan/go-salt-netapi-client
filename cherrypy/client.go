// Package cherrypy provides a client to integrate with Salt NetAPI's rest_cherrypy module
// https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html
package cherrypy

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type eauth struct {
	Username string
	Password string
	Backend  string
}

/*
Client handles communication with NetAPI rest_cherrypy module (https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html)

Example usage:
	client := cherrypy.NewClient("http://master:8000", "admin", "password", "pam")
	if err := client.Login(); err != nil {
		return err
	}
	defer client.Logout()

	minion := client.Minion("minion1")
*/
type Client struct {
	client  *http.Client
	eauth   *eauth
	Address string
	Token   string
}

/*
NewClient creates a new instance of client
  address: URL of the cherrypy instance on a master (e.g.: https://salt-master:8000)
  backend: External authentication (eauth) backend (https://docs.saltstack.com/en/latest/topics/eauth/index.html)
*/
func NewClient(address string, username string, password string, backend string) *Client {
	a := eauth{
		Username: username,
		Password: password,
		Backend:  backend,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &Client{
		client:  &http.Client{Transport: tr},
		eauth:   &a,
		Address: address,
	}
}

func (c *Client) request(method string, endpoint string, accept string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.Address, endpoint)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Sending request %s to %s", jsonData, url)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", accept)
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("X-Auth-Token", c.Token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Received response (%d) %s from %s", resp.StatusCode, body, url)
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("HTTP Request failed: %s.\n%s", resp.Status, body)
	}

	return body, nil
}

func (c *Client) requestJSON(method string, endpoint string, data interface{}) (map[string]interface{}, error) {
	body, err := c.request(method, endpoint, "application/json", data)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
