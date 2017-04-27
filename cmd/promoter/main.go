package main

import (
	"encoding/json"
	"net/http"

	"github.com/unrolled/render"
)

func main() {
}

type promotion struct {
	ID   string `json:"id"`
	Snap string `json:"snap"`
}

type Persister interface {
	AddPromotion(*promotion) string
	GetPromotion(string) *promotion
}

func createPromotionHandler(formatter *render.Render, repo Persister) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		dec := json.NewDecoder(req.Body)
		p := promotion{}
		dec.Decode(&p)
		id := repo.AddPromotion(&p)
		w.Header().Add("Location", "/promotions/1")
		formatter.JSON(w,
			http.StatusCreated,
			&promotion{
				ID: id,
			},
		)
	}
}
