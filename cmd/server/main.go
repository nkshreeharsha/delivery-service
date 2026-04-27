package main

import (
	"log"
	"net/http"
	"time"
	"github.com/nkshreeharsha/delivery-service/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	signer := &stubSigner{url: "https://example.com/signed-url"}

	h := handler.New(signer,"2026-04-22_1",60*time.Second)

	r := chi.NewRouter()
	r.Use(middleware.Logger) 

	r.Get("/v1/subscriber/{subID}/creativeList", h.GetCreativeList)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Placeholder — swap for real GCS implementation
type stubSigner struct{ url string}

func (s *stubSigner) SignURL(p string, _ time.Duration) (string, error) {
    return "https://storage.googleapis.com/fake/" + p, nil
}

func newGCSSigner() handler.StorageSigner { return &stubSigner{} }