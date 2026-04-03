package location

import (
	"net/http"

	"github.com/femitubosun/go-sweepline-availability/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}

func (h *handler) GetLocationInfo(w http.ResponseWriter, r *http.Request) {
	defaultLocation, err := h.service.GetInfo(r.Context())

	if err != nil {
		json.Write(w, http.StatusInternalServerError, nil)
		return
	}

	json.Write(w, http.StatusOK, toLocationResponse(defaultLocation))
}
