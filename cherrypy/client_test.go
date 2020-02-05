package cherrypy

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	apiTester "github.com/finarfin/go-apiclient-tester/common"
	"github.com/finarfin/go-apiclient-tester/postman"
)

const (
	testUsername    = "test_user"
	testPassword    = "test_pwd"
	testEAuth       = "pam"
	testToken       = "163588fd62e0166d48196be8dbfec35287931f10"
	testSampleJobID = "20200202210231414902"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func setup(t *testing.T) (*apiTester.Tester, *Client) {
	tester, err := postman.NewTester("testdata/cherrypi_collection.json")
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient(tester.URL, testUsername, testPassword, testEAuth)
	client.Token = testToken

	return tester, client
}
