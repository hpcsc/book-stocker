package info

import (
	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"
	"net/http"
)

var (
	Version = "main"
)

type infoResponse struct {
	Version string
}

type handler struct {
	renderer *render.Render
}

func RegisterRoutes(router *chi.Mux) {
	h := handler{renderer: render.New()}
	router.Get("/info", h.get)
}

func (h *handler) get(w http.ResponseWriter, req *http.Request) {
	_ = h.renderer.JSON(w, http.StatusOK, infoResponse{Version: Version})
}
