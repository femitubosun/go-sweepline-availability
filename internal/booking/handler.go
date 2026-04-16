package booking

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	appJSON "github.com/femitubosun/go-sweepline-availability/internal/json"
)

type handler struct {
	service Service
}

func (h *handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	courtID := r.PathValue("courtID")

	if courtID == "" {
		appJSON.Write(w, http.StatusBadRequest, map[string]string{"error": "courtID is required"})
		return
	}

	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appJSON.Write(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	input, err := req.toCreateBookingInput(courtID)
	if err != nil {
		appJSON.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	booking, err := h.service.CreateBooking(ctx, input)
	if err != nil {
		if errors.Is(err, ErrHistoricalBooking) {
			appJSON.Write(w, http.StatusBadRequest, map[string]string{"error": "booking cannot be in the past"})
			return
		}
		appJSON.Write(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	appJSON.Write(w, http.StatusCreated, toBookingResponse(booking))
}

func (h *handler) GetCourtBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	courtID := r.PathValue("courtID")

	bookings, err := h.service.ListCourtBookings(ctx, courtID)
	if err != nil {
		appJSON.Write(w, http.StatusInternalServerError, map[string]string{"error": "failed to get bookings"})
		return
	}

	resp := make([]bookingResponse, len(bookings))
	for i, b := range bookings {
		resp[i] = toBookingResponse(b)
	}

	appJSON.Write(w, http.StatusOK, resp)
}

func (h *handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bookingID := r.PathValue("bookingID")

	err := h.service.CancelBooking(ctx, bookingID)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			appJSON.Write(w, http.StatusNotFound, map[string]string{"error": "booking not found"})
			return
		}
		appJSON.Write(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	appJSON.Write(w, http.StatusNoContent, nil)
}

func (h *handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bookingID := r.PathValue("bookingID")

	booking, err := h.service.GetBooking(ctx, bookingID)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			appJSON.Write(w, http.StatusNotFound, map[string]string{"error": "booking not found"})
			return
		}
		appJSON.Write(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	appJSON.Write(w, http.StatusOK, toBookingResponse(booking))
}

func (h *handler) AvailabilitySearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	courtID := r.PathValue("courtID")

	fmt.Print("Availability search")

	from, to, durationMins, err := parseAvailabilityQuery(r)
	if err != nil {
		appJSON.Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if from.Before(time.Now()) {
		appJSON.Write(w, http.StatusBadRequest, map[string]string{"error": "cannot search for availability in the past"})
		return
	}

	gaps, err := h.service.GetGaps(ctx, GetGapsInput{
		CourtID:      courtID,
		DurationMins: durationMins,
		Interval:     Interval{Start: from, End: to},
	})
	if err != nil {
		appJSON.Write(w, http.StatusInternalServerError, map[string]string{"error": "failed to search availability"})
		return
	}

	resp := make([]availabilitySlotResponse, len(gaps))
	for i, gap := range gaps {
		resp[i] = availabilitySlotResponse{
			Start: gap.Start.Format(time.RFC3339),
			End:   gap.End.Format(time.RFC3339),
		}
	}

	appJSON.Write(w, http.StatusOK, resp)
}

func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}

type createBookingRequest struct {
	StartTime       string `json:"startTime"`
	DurationMinutes int    `json:"durationMinutes"`
	GuestName       string `json:"guestName"`
}

func (r createBookingRequest) toCreateBookingInput(courtID string) (CreateBookingInput, error) {
	start, err := time.Parse(time.RFC3339, r.StartTime)
	if err != nil {
		return CreateBookingInput{}, fmt.Errorf("startTime: must be RFC3339")
	}

	return CreateBookingInput{
		CourtID:         courtID,
		StartTime:       start,
		DurationMinutes: r.DurationMinutes,
		GuestName:       r.GuestName,
	}, nil
}

type bookingResponse struct {
	ID        string `json:"id"`
	CourtID   string `json:"courtId"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	GuestName string `json:"guestName"`
}

func toBookingResponse(b Booking) bookingResponse {
	return bookingResponse{
		ID:        b.ID,
		CourtID:   b.CourtID,
		StartTime: b.Interval.Start.Format(time.RFC3339),
		EndTime:   b.Interval.End.Format(time.RFC3339),
		GuestName: b.GuestName,
	}
}

type availabilitySlotResponse struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func parseAvailabilityQuery(r *http.Request) (from, to time.Time, durationMins int, err error) {
	q := r.URL.Query()

	fromStr := q.Get("from")
	if fromStr == "" {
		return from, to, 0, fmt.Errorf("from: required")
	}
	from, err = time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return from, to, 0, fmt.Errorf("from: must be RFC3339")
	}

	toStr := q.Get("to")
	if toStr == "" {
		return from, to, 0, fmt.Errorf("to: required")
	}
	to, err = time.Parse(time.RFC3339, toStr)
	if err != nil {
		return from, to, 0, fmt.Errorf("to: must be RFC3339")
	}

	durationStr := q.Get("duration_minutes")
	if durationStr == "" {
		return from, to, 0, fmt.Errorf("duration_minutes: required")
	}
	durationMins, err = strconv.Atoi(durationStr)
	if err != nil {
		return from, to, 0, fmt.Errorf("duration_minutes: must be integer")
	}
	if durationMins <= 0 {
		return from, to, 0, fmt.Errorf("duration_minutes: must be positive")
	}

	return from, to, durationMins, nil
}
