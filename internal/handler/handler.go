package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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


func (h *Handler) GetCreativeList(c *gin.Context){
	subID := c.Param("subID")
	if subID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error" : "0x01"})
		return
	}

	objectPath := fmt.Sprintf("%s/%s.csv", h.activeFolderVersion, subID)
	log.Printf("Generating signed URL for object path: %s", objectPath)
	signedURL, err := h.signer.SignURL(objectPath, h.signedURLTTL)
	log.Printf("Signed URL generation result: %s, error: %v", signedURL, err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "0x02"})
		return
	}

	c.Redirect(http.StatusFound, signedURL)
}