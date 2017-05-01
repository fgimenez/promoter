package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

type fakePersister struct{}

var insertedPromotions []*promotion
var currentID int

func (f *fakePersister) AddPromotion(p *promotion) string {
	currentID++
	p.ID = strconv.Itoa(currentID)
	if insertedPromotions == nil {
		insertedPromotions = []*promotion{p}
	} else {
		insertedPromotions = append(insertedPromotions, p)
	}
	return p.ID
}

func (f *fakePersister) GetPromotion(id string) *promotion {
	for _, p := range insertedPromotions {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func TestCreatePromotion(t *testing.T) {
	client := &http.Client{}

	repo := &fakePersister{}

	server := httptest.NewServer(
		http.HandlerFunc(createPromotionHandler(formatter, repo)))
	defer server.Close()

	values := map[string]string{
		"snap":          "core",
		"architecture":  "amd64",
		"revision":      "1931",
		"status":        "passed",
		"last_update":   "2017-04-06T13:42:31Z",
		"signed_off_by": "fgimenez",
		"comments":      "There are some failing tests that will be fixed when snapd#3018 lands in the release branch",
	}

	body, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer([]byte(body)))

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
	t.Run("persistence", func(t *testing.T) {
		t.Run("first insertion", func(t *testing.T) {
			insertedPromotion := repo.GetPromotion(p.ID)
			if insertedPromotion.Snap != "core" {
				t.Error("inserted promotion not found")
			}
		})
		t.Run("additional insertion", func(t *testing.T) {
			body := []byte(`{
                    "snap": "core2",
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

			dec := json.NewDecoder(resp.Body)
			p := promotion{}
			if err := dec.Decode(&p); err != nil {
				t.Errorf("error decoding response %v", err)
			}
			insertedPromotion := repo.GetPromotion(p.ID)
			if insertedPromotion.Snap != "core2" {
				t.Error("inserted promotion not found")
			}
		})
		t.Run("architecture field", func(t *testing.T) {
			expected := values["architecture"]
			if p.Architecture != expected {
				t.Errorf("field in promotion with wrong value, expected %s, found, %s", expected, p.Architecture)
			}
		})
		t.Run("revision field", func(t *testing.T) {
			expected := values["revision"]
			if p.Revision != expected {
				t.Errorf("field in promotion with wrong value, expected %s, found, %s", expected, p.Revision)
			}
		})
		t.Run("status field", func(t *testing.T) {
			expected := values["status"]
			if p.Status != expected {
				t.Errorf("field in promotion with wrong value, expected %s, found, %s", expected, p.Status)
			}
		})
		t.Run("last_update field", func(t *testing.T) {
			expected := values["last_update"]
			if p.LastUpdate != expected {
				t.Errorf("field in promotion with wrong value, expected %s, found, %s", expected, p.LastUpdate)
			}
		})
		t.Run("siged_off_by field", func(t *testing.T) {
			expected := values["signed_off_by"]
			if p.SignedOffBy != expected {
				t.Errorf("field in promotion with wrong value, expected %s, found, %s", expected, p.SignedOffBy)
			}
		})
		t.Run("comments field", func(t *testing.T) {
			expected := values["comments"]
			if p.Comments != expected {
				t.Errorf("field in promotion with wrong value, expected %s, found, %s", expected, p.Comments)
			}
		})

	})
}
