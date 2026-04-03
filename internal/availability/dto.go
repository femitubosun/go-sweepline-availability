package availability

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type GetAvailabilityResponse struct {
	CourtID        string `json:"courtId"`
	AvailableSlots AvailabilitySlotsResponse
}

type AvailabilitySlotsResponse struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type GetAvailabilityQuery struct {
	CourtID         string
	From            time.Time `query:"from" validate:"required"`
	To              time.Time `query:"to" validate:"required,gtfield=From"`
	DurationMinutes int       `query:"duration_minutes" validate:"required,min=1"`
}

func parseGetAvailabilityRequest(r *http.Request) (GetAvailabilityQuery, error) {
	q := r.URL.Query()
	var req GetAvailabilityQuery

	if fromStr := q.Get("from"); strings.TrimSpace(fromStr) != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return req, fmt.Errorf("from: must be RFC3339")
		}
		req.From = t
	}

	if toStr := q.Get("to"); strings.TrimSpace(toStr) != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return req, fmt.Errorf("to: must be RFC3339")
		}
		req.To = t
	}

	if durStr := q.Get("duration_minutes"); strings.TrimSpace(durStr) != "" {
		d, err := strconv.Atoi(durStr)
		if err != nil {
			return req, fmt.Errorf("duration_minutes: must be integer")
		}
		req.DurationMinutes = d
	}

	if err := validator.New().Struct(req); err != nil {
		return req, fmt.Errorf("validation failed: %w", err)
	}

	return req, nil
}
