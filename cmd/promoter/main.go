package main

import (
	"net/http"

	"github.com/unrolled/render"
)

func main() {
}

func createPromotionHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Location", "/promotions/")
		formatter.JSON(w,
			http.StatusCreated,
			struct{ Test string }{"This is a test"})
	}
}
