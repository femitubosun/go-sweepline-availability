package booking

import "net/http"

type handler struct {
	service Service
}

func (h *handler) CreateBooking(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) CancelBooking(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) GetCourtBookings(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) AvailabilitySearch(w http.ResponseWriter, r *http.Request) {

}

func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}
