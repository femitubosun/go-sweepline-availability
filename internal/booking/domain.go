package booking

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusActive    BookingStatus = "active"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID              string        `json:"id"`
	GuestName       string        `json:"guestName"`
	CourtID         string        `json:"courtId"`
	Interval        Interval      `json:"interval"`
	DurationMinutes int           `json:"durationMinutes"`
	Status          BookingStatus `json:"status"`
	CancelledAt     *time.Time    `json:"cancelledAt,omitempty"`
}
