package availability

import (
	"net/http"

	"github.com/femitubosun/go-sweepline-availability/internal/json"
)

type handler struct {
	service Service
}

func (h *handler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	req, err := parseGetAvailabilityRequest(r)
	if err != nil {
		json.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	result, err := h.service.GetAvailability(ctx, GetAvailabilityQuery{
		CourtID:         id,
		DurationMinutes: req.DurationMinutes,
		From:            req.From,
		To:              req.To,
	})

	if err != nil {
		json.Write(w, http.StatusInternalServerError, nil)
		return
	}

	json.Write(w, http.StatusOK, result)

}

func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}
