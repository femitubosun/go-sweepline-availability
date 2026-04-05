package booking

type CreateBookingRequest struct {
	CourtID         string `json:"courtId" validate:"required"`
	From            string `json:"from" validate:"required"`
	DurationMinutes int    `json:"durationMinutes" validate:"required,min=1"`
	GuestName       string `json:"guestName" validate:"required"`
}

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "Pending"
	BookingStatusActive    BookingStatus = "Active"
	BookingStatusCompleted BookingStatus = "Completed"
	BookingStatusCancelled BookingStatus = "Cancelled"
)

type BookingResponse struct {
	ID              string        `json:"id"`
	GuestName       string        `json:"guestName"`
	CourtID         string        `json:"courtId"`
	StartTime       string        `json:"startTime"`
	EndTime         string        `json:"endTime"`
	DurationMinutes int           `json:"durationMinutes"`
	Status          BookingStatus `json:"status"`
	CancelledAt     string        `json:"cancelledAt"`
}
