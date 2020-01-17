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
func NewClient(address string, username string, password string, eauth string) (*Client, error) {
	a := EAuth{}
	a.Username = username
	a.Password = password
	a.Backend = eauth

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	c := Client{
		Client:  &http.Client{Transport: tr},
		Address: address,
		EAuth:   &a,
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Client) sendRequest(method string, endpoint string, data interface{}) (map[string]interface{}, error) {
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

	log.Printf("[DEBUG] Received response %s from %s", body, url)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[WARN] HTTP Request failed: %s.\n%s", resp.Status, body)
	}

	if resp.Header.Get("Content-Type") == "application/json" {
		result := make(map[string]interface{})
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, nil
}
