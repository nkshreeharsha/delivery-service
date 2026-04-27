package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

)

type StorageSigner interface {
	SignURL(path string, ttl time.Duration) (string, error)
}

type Handler struct {
	signer StorageSigner
	activeFolderVersion string
	signedURLTTL time.Duration

}

func New(signer StorageSigner, activeFolderVersion string, signedURLTTL time.Duration) *Handler {
	return &Handler{
		signer: signer, 
		activeFolderVersion: activeFolderVersion,
		signedURLTTL: signedURLTTL}
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, errorCode string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	http.Error(w, fmt.Sprintf(`{"error": "%s"}`, errorCode), statusCode)
	fmt.Fprint(w, message)
}

func (h *Handler) GetCreativeList(w http.ResponseWriter, r *http.Request){
	subID := chi.URLParam(r, "subID")
	if subID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "0x01", "Missing obfuscatedSubId")
		return
	}

	// Placeholder for actual creative list retrieval logic
	objectPath := fmt.Sprintf("%s/%s.csv", h.activeFolderVersion, subID)

	signedURL, err := h.signer.SignURL(objectPath, h.signedURLTTL)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "0x02", "Failed to sign URL")
		return
	}

	http.Redirect(w, r, signedURL, http.StatusFound)

}