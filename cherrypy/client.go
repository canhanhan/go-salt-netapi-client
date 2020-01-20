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

// EAuth contains authenticated data
type EAuth struct {
	Username string
	Password string
	Backend  string
}

// Client handles communication with CherryPy
type Client struct {
	Client  *http.Client
	EAuth   *EAuth
	Address string
	Token   string
}

// NewClient creates a new instance of client
func NewClient(address string, username string, password string, eauth string) *Client {
	a := EAuth{
		Username: username,
		Password: password,
		Backend:  eauth,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &Client{
		Client:  &http.Client{Transport: tr},
		Address: address,
		EAuth:   &a,
	}
}

func (c *Client) request(method string, endpoint string, data interface{}) ([]byte, error) {
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

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("X-Auth-Token", c.Token)
	}

	resp, err := c.Client.Do(req)
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
	body, err := c.request(method, endpoint, data)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
