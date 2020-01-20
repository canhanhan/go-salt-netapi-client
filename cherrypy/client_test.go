package cherrypy

import (	
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	testUsername = "sample_user"
	testPassword = "sample_password"
	testEAuth = "file"
	testToken = "75023210fea33137fd41d24d75998b93eba9b103"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func setup(t *testing.T) (*Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(server.URL, testUsername, testPassword, testEAuth)

	return client, mux, server.Close
}

func handleJSONRequest(mux *http.ServeMux, path string, scenario string) {
	handleRequestWithHeaders(mux, path, scenario, map[string]string{
		"Content-Type": "application/json",
		"X-Auth-Token": testToken,
	})
}

func handleRequest(mux *http.ServeMux, path string, scenario string) {
	handleRequestWithHeaders(mux, path, scenario, map[string]string{})
}

func handleRequestWithHeaders(mux *http.ServeMux, path string, scenario string, headers map[string]string) {
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		content, err := getResponseText(scenario)
		if err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			for k, v := range headers {
				w.Header().Set(k, v)
			}
			
			fmt.Fprintf(w, content)
		}	
	})
}

func getResponseText(name string) (string, error) {
	data, err := getResponse(fmt.Sprintf("%s.json", name))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getResponse(name string) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("testdata/response/%s", name))
}