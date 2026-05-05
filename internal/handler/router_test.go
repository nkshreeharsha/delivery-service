package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	
    "github.com/gin-gonic/gin"
	"github.com/nkshreeharsha/delivery-service/internal/handler"
	"github.com/stretchr/testify/assert"
)

func buildRouter(signer *fakeSigner) http.Handler {
	gin.SetMode(gin.TestMode)
    h := handler.New(signer, "2026-04-22_1", 60*time.Second)
    r := gin.New()
    r.GET("/v1/subscriber/:subID/creativeList", h.GetCreativeList)
    return r
}

func Test_router_validPath_returns302(t *testing.T) {
	signer := &fakeSigner{returnURL: "https://aws-s3-signed-url.com"}
	srv := httptest.NewServer(buildRouter(signer))
	defer srv.Close()

	client := &http.Client{
        CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
            return http.ErrUseLastResponse
        },
    }

	response , err := client.Get(srv.URL + "/v1/subscriber/abc123/creativeList")
	assert.NoError(t, err)
	defer response.Body.Close()

	assert.Equal(t, http.StatusFound, response.StatusCode)
	assert.Equal(t, signer.returnURL, response.Header.Get("Location"))
}

func Test_router_invalidPath_returns404(t *testing.T) {
	srv := httptest.NewServer(buildRouter(&fakeSigner{}))
	defer srv.Close()

	response,_ := http.Get(srv.URL + "/invalid/path")
	assert.NotNil(t, response)
	defer response.Body.Close()
	

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}
