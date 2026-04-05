package booking

type CreateBookingRequest struct {
	CourtID         string `json:"courtId" validate:"required"`
	From            string `json:"from" validate:"required"`
	DurationMinutes int    `json:"durationMinutes" validate:"required,min=1"`
	GuestName       string `json:"guestName" validate:"required"`
}
