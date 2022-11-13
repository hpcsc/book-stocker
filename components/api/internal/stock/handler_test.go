//go:build unit

package stock

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/hpcsc/book-stocker/api/internal/store"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Post(t *testing.T) {
	t.Run("return bad request when request is not valid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		router, _ := testRouterWithStockRoutes()
		httpRequest, err := http.NewRequest("POST", "/stock", strings.NewReader("test: json"))
		require.NoError(t, err)
		router.ServeHTTP(w, httpRequest)

		require.Equal(t, http.StatusBadRequest, w.Code)
		response := decodeStockResponseBody(t, w.Body)
		require.False(t, response.Successful)
		require.Contains(t, response.Error, "failed to unmarshal stock request")
	})

	t.Run("return bad request when isbn is not provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		router, _ := testRouterWithStockRoutes()
		httpRequest := newStockHttpRequest(t, stockRequest{
			Quantity: 1,
		})
		router.ServeHTTP(w, httpRequest)

		require.Equal(t, http.StatusBadRequest, w.Code)
		response := decodeStockResponseBody(t, w.Body)
		require.False(t, response.Successful)
		require.Equal(t, "ISBN is required", response.Error)
	})

	t.Run("return bad request when quantity is not provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		router, _ := testRouterWithStockRoutes()
		httpRequest := newStockHttpRequest(t, stockRequest{
			ISBN: "some-isbn",
		})
		router.ServeHTTP(w, httpRequest)

		require.Equal(t, http.StatusBadRequest, w.Code)
		response := decodeStockResponseBody(t, w.Body)
		require.False(t, response.Successful)
		require.Equal(t, "Quantity must be greater than 0", response.Error)
	})

	t.Run("return accepted with id when stock request is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		router, requestStore := testRouterWithStockRoutes()
		requestStore.StubSave().Return(nil)
		httpRequest := newStockHttpRequest(t, stockRequest{
			ISBN:     "some-isbn",
			Quantity: 2,
		})
		router.ServeHTTP(w, httpRequest)

		require.Equal(t, http.StatusAccepted, w.Code)
		response := decodeStockResponseBody(t, w.Body)
		require.True(t, response.Successful)
		require.Empty(t, response.Error)
		require.NotEmpty(t, response.Id)
	})

	t.Run("return internal error when fail to save to store", func(t *testing.T) {
		w := httptest.NewRecorder()
		router, requestStore := testRouterWithStockRoutes()
		requestStore.StubSave().Return(errors.New("some error"))
		httpRequest := newStockHttpRequest(t, stockRequest{
			ISBN:     "some-isbn",
			Quantity: 2,
		})
		router.ServeHTTP(w, httpRequest)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		response := decodeStockResponseBody(t, w.Body)
		require.False(t, response.Successful)
		require.Contains(t, response.Error, "some error")
	})
}

func testRouterWithStockRoutes() (*chi.Mux, *store.Fake) {
	router := chi.NewRouter()
	requestStore := store.NewFake()
	v := validator.New()

	RegisterRoutes(router, v, requestStore)

	return router, requestStore
}

func newStockHttpRequest(t *testing.T, request stockRequest) *http.Request {
	marshalledRequest, err := json.Marshal(request)
	httpRequest, err := http.NewRequest("POST", "/stock", bytes.NewBuffer(marshalledRequest))
	require.NoError(t, err)
	return httpRequest
}

func decodeStockResponseBody(t *testing.T, body *bytes.Buffer) stockResponse {
	var response stockResponse
	err := json.NewDecoder(body).Decode(&response)
	require.NoError(t, err)
	return response
}
