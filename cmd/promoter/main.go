package main

import (
	"net/http"

	"github.com/unrolled/render"
)

func main() {
}

type promotion struct {
	ID string `json:"id"`
}

func createPromotionHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Location", "/promotions/1")
		formatter.JSON(w,
			http.StatusCreated,
			&promotion{ID: "1"},
		)
	}
}
