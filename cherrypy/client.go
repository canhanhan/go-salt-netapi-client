// Package cherrypy provides a client to integrate with Salt NetAPI's rest_cherrypy module
// https://docs.saltstack.com/en/latest/ref/netapi/all/salt.netapi.rest_cherrypy.html
package cherrypy

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type eauth struct {
	Username string
	Password string
	Backend  string
}

type saltUnixTime struct {
	time.Time
}

func (t *saltUnixTime) UnmarshalJSON(input []byte) error {
	s, err := strconv.ParseFloat(string(input), 64)
	if err != nil {
		return err
	}

	m := int64(s)
	n := int64((s - float64(m)) * 1000000000)
	t.Time = time.Unix(m, n)
	return nil
}

type saltTime struct {
	time.Time
}

func (t *saltTime) UnmarshalJSON(input []byte) error {
	s := string(input)
	s = strings.Trim(s, "\"")
	v, err := time.Parse("2006, Jan 02 15:04:05.000000", s)
	if err != nil {
		return err
	}

	t.Time = v
	return nil
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

func (c *Client) newRequest(ctx context.Context, method string, endpoint string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.Address, endpoint)

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("[DEBUG] Creating request for %s", url)
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("X-Auth-Token", c.Token)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	
	log.Printf("[DEBUG] Received response %s from %s", resp.Status, resp.Request.URL)
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("HTTP Request failed: %s", resp.Status)
	}
	
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil && err != io.EOF {
				return nil, err
			}
		}
	}

	return resp, nil
}
