package main

import (
	"bytes"
	"encoding/json"
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
		http.HandlerFunc(createPromotionHandler(formatter)))
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
		t.Errorf("Errored when sending request to the server %v", err)
	}
	defer resp.Body.Close()

	t.Run("statusCode", func(t *testing.T) {
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected response status 201, received %s",
				resp.Status)
		}
	})
	dec := json.NewDecoder(resp.Body)
	p := promotion{}
	if err := dec.Decode(&p); err != nil {
		t.Errorf("error decoding response %v", err)
	}

	t.Run("locationHeader", func(t *testing.T) {
		loc, headerOk := resp.Header["Location"]
		t.Run("Set", func(t *testing.T) {
			if !headerOk {
				t.Error("Location header is not set")
			}
		})
		t.Run("Path", func(t *testing.T) {
			if !strings.Contains(loc[0], "/promotions/") {
				t.Errorf("Location header should contain '/promotions/'")
			}
		})
		t.Run("ID", func(t *testing.T) {
			locationItems := strings.Split(loc[0], "/")
			if p.ID == "" {
				t.Error("empty p.ID in payload")
			}
			if p.ID != locationItems[2] {
				t.Errorf("id from payload %v, from location %v", p.ID, locationItems[2])
			}
		})
	})
}
