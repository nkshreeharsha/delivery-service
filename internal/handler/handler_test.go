package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nkshreeharsha/delivery-service/internal/handler"

	"github.com/go-chi/chi/v5"
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

func requestWithValidSubID(subID string) *http.Request {
    req := httptest.NewRequest(http.MethodGet, "/v1/subscriber/"+subID+"/creativeList", nil)
    rctx := chi.NewRouteContext()
    rctx.URLParams.Add("subID", subID)
    return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func Test_emptySubID_returns400(t *testing.T) {
	h := handler.New(&fakeSigner{}, "2026-04-22_1", 60*time.Second)

	req := httptest.NewRequest(http.MethodGet, "/any", nil)
	w := httptest.NewRecorder()

	h.GetCreativeList(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "0x01")
}

func Test_validSubId_returns302(t *testing.T) {
    expectedURL := "https://example.com/signed-url"
    fakeSigner := &fakeSigner{returnURL: expectedURL}
    h := handler.New(fakeSigner, "2026-04-22_1", 60*time.Second)

    w := httptest.NewRecorder()
    h.GetCreativeList(w, requestWithValidSubID("abc123"))

    assert.Equal(t, http.StatusFound, w.Code)  // 302
    assert.Equal(t, expectedURL, w.Header().Get("Location"))
   
}

func Test_validSubId_signerError_returns500(t *testing.T) {
    fakeSigner := &fakeSigner{returnError: errors.New("Signing failed")}
    h := handler.New(fakeSigner, "2026-04-22_1", 60*time.Second)

    w := httptest.NewRecorder()
    h.GetCreativeList(w, requestWithValidSubID("abc123"))

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    assert.Contains(t, w.Body.String(), "0x02")
}