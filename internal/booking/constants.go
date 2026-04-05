package booking

const (
	MaxBookingDurationDays = 3

	RedisKeyBookings    = "bookings:%s"
	RedisKeyBookingEnds = "booking_ends"
)
