package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hpcsc/book-stocker/api/internal/store"
	"github.com/hpcsc/book-stocker/api/internal/validate"
	"github.com/unrolled/render"
	"net/http"
)

type stockRequest struct {
	ISBN     string `json:"isbn" validate:"required"`
	Quantity int    `json:"quantity" validate:"gt=0"`
}

type stockResponse struct {
	Id         string `json:"id"`
	Successful bool   `json:"successful"`
	Error      string `json:"error"`
}

type handler struct {
	renderer     *render.Render
	validator    *validator.Validate
	requestStore store.Interface
}

func RegisterRoutes(router *chi.Mux, validator *validator.Validate, requestStore store.Interface) {
	h := &handler{
		renderer:     render.New(),
		validator:    validator,
		requestStore: requestStore,
	}
	router.Post("/stock", h.post)
}

func (h *handler) post(w http.ResponseWriter, req *http.Request) {
	var request stockRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		h.renderBadRequestStockResponse(w, fmt.Sprintf("failed to unmarshal stock request: %v", err))
		return
	}

	if err := h.validator.Struct(request); err != nil {
		h.renderBadRequestStockResponse(w, validate.Reformat(err).Error())
		return
	}

	id := uuid.New().String()
	if err := h.requestStore.Save(context.TODO(), store.StockRequest{
		Id:       id,
		ISBN:     request.ISBN,
		Quantity: request.Quantity,
	}); err != nil {
		_ = h.renderer.JSON(w, http.StatusInternalServerError, stockResponse{
			Successful: false,
			Error:      err.Error(),
		})
		return
	}

	_ = h.renderer.JSON(w, http.StatusAccepted, stockResponse{
		Successful: true,
		Id:         id,
	})
}

func (h *handler) renderBadRequestStockResponse(w http.ResponseWriter, message string) {
	_ = h.renderer.JSON(w, http.StatusBadRequest, stockResponse{
		Successful: false,
		Error:      message,
	})
}
