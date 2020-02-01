package cherrypy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testUsername = "test_username"
	testPassword = "test"
	testEAuth    = "pam"
	testToken    = "eaffa1fd47cf4254dbcbcf2c0145f8b5f78d1e70"
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

type responseReply struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func parseFile(name string) (hdr string, body string, err error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s", name))
	if err != nil {
		return
	}

	var hEnd, bStart int
	for i, v := range b {
		if i+1 < len(b) && v == '\n' && b[i+1] == '\n' {
			hEnd = i
			bStart = i + 2
			break
		} else if i+3 < len(b) && v == '\r' && b[i+1] == '\n' && b[i+2] == '\r' && b[i+3] == '\n' {
			hEnd = i
			bStart = i + 4
			break
		}
	}

	if bStart == 0 || hEnd == 0 {
		err = fmt.Errorf("Could not parse the file. hEnd: %d, bStart: %d", hEnd, bStart)
		return
	}

	hdr = string(b[0:hEnd])
	body = string(b[bStart:])

	return
}

func getContentFromFile(name string) (*strings.Reader, error) {
	hdr, body, err := parseFile(name)
	if err != nil {
		return nil, err
	}

	hdr += fmt.Sprintf("\nContent-Length: %d", len(body))

	return strings.NewReader(hdr + "\n\n" + body), nil
}

func setupScenario(t *testing.T, mux *http.ServeMux, scenario string) {
	reqf, err := getContentFromFile(scenario + ".req")
	if err != nil {
		t.Fatal(err)
	}

	request, err := http.ReadRequest(bufio.NewReader(reqf))
	if err != nil {
		t.Fatal(err)
	}

	expectedRequestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		t.Fatal(err)
	}

	path := request.URL.Path

	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		f, err := getContentFromFile(scenario + ".res")
		if err != nil {
			t.Fatal(err)
		}

		actualRequestBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}

		response, err := http.ReadResponse(bufio.NewReader(f), request)
		if err != nil {
			t.Fatal(err)
		}
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		// Compare request
		assert.Equal(t, request.Method, req.Method)
		for k, v := range request.Header {
			if k == "Content-Length" && request.Header.Get("Content-Type") == "application/json" {
				continue
			}

			if value, ok := req.Header[k]; !ok {
				assert.Failf(t, "Request headers do not contain %s", k)
			} else {
				assert.Equal(t, v, value)
			}
		}

		if request.Header.Get("Content-Type") == "application/json" {
			expectedRequestData := make(map[string]interface{})
			if err := json.Unmarshal(expectedRequestBody, &expectedRequestData); err != nil {
				t.Fatal(err)
			}

			actualRequestData := make(map[string]interface{})
			if err := json.Unmarshal(actualRequestBody, &actualRequestData); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expectedRequestData, actualRequestData)
		} else {
			assert.Equal(t, expectedRequestBody, actualRequestBody)
		}

		// Send response
		w.WriteHeader(response.StatusCode)
		for k, vs := range response.Header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}

		_, err = w.Write(responseBody)
		if err != nil {
			t.Fatal(err)
		}
	})
}
