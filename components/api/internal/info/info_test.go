//go:build unit

package info

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInfoHandler_Get(t *testing.T) {
	t.Run("return current version", func(t *testing.T) {
		w := httptest.NewRecorder()
		router := testRouterWithInfoRoutes()
		httpRequest, err := http.NewRequest("GET", "/info", nil)
		require.NoError(t, err)
		router.ServeHTTP(w, httpRequest)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "main")
	})
}

func testRouterWithInfoRoutes() *chi.Mux {
	router := chi.NewRouter()

	RegisterRoutes(router)

	return router
}
