package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func TestCreatePromotion(t *testing.T) {
	client := &http.Client{}

	server := httptest.NewServer(
		http.HandlerFunc(createMatchHandler(formatter)))
	defer server.Close()

	body := []byte(`{
                    "snap": "core",
                    "architecture": "amd64",
                    "revision": 1931,
                    "status": "passed",
                    "last_update": "2017-04-06T13:42:31Z ",
                    "signed_off_by": "fgimenez",
                    "comments": "There are some failing tests that will be fixed when snapd#3018 lands in the release branch"
                  }`)

	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status 201, received %s",
			resp.Status)
	}

	loc, headerOk := resp.Header["Location"]
	if !headerOk {
		t.Error("Location header is not set")
	}

	if !strings.Contains(loc[0], "/promotions/") {
		t.Errorf("Location header should contain '/promotions/'")
	}

	fmt.Printf("Payload: %s", string(payload))
}
