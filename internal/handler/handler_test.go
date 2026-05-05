package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nkshreeharsha/delivery-service/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeSigner struct{
    returnURL string
    returnError error
    gotPath string
}

func (f *fakeSigner) SignURL(path string, _ time.Duration) (string, error) {
	f.gotPath = path
    return f.returnURL,f.returnError
}

func Test_emptySubID_returns400(t *testing.T) {
	h := handler.New(&fakeSigner{}, "2026-04-22_1", 60*time.Second)
	w := httptest.NewRecorder()
    c,_ := gin.CreateTestContext(w)

	h.GetCreativeList(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "0x01")
}

func Test_validSubId_returns302(t *testing.T) {
    expectedURL := "https://example.com/signed-url"
    fakeSigner := &fakeSigner{returnURL: expectedURL}
    h := handler.New(fakeSigner, "2026-04-22_1", 60*time.Second)

    w := httptest.NewRecorder()
    c,_ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodGet, "/v1/subscriber/abc123/creativeList", nil)
    c.Params = gin.Params{{Key: "subID", Value: "abc123"}}
    h.GetCreativeList(c)

    assert.Equal(t, http.StatusFound, w.Code)  // 302
    assert.Equal(t, expectedURL, w.Header().Get("Location"))
   
}

func Test_validSubId_signerError_returns500(t *testing.T) {
    fakeSigner := &fakeSigner{returnError: errors.New("Signing failed")}
    h := handler.New(fakeSigner, "2026-04-22_1", 60*time.Second)

    w := httptest.NewRecorder()
    c,_ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest(http.MethodGet, "/v1/subscriber/abc123/creativeList", nil)
    c.Params = gin.Params{{Key: "subID", Value: "abc123"}}
    h.GetCreativeList(c)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    assert.Contains(t, w.Body.String(), "0x02")
}