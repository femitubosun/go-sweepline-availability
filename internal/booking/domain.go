package booking

type Booking struct {
	ID              string        `json:"id"`
	GuestName       string        `json:"guestName"`
	CourtID         string        `json:"courtId"`
	StartTime       string        `json:"startTime"`
	EndTime         string        `json:"endTime"`
	DurationMinutes int           `json:"durationMinutes"`
	Status          BookingStatus `json:"status"`
	CancelledAt     string        `json:"cancelledAt"`
}
