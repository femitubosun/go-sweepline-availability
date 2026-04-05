package booking

import (
	"github.com/google/uuid"
)

func generateBookingId() string {
	return uuid.New().String()
}
