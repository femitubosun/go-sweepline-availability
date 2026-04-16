package main

import (
	"net/http"

	"github.com/femitubosun/go-sweepline-availability/internal/booking"
	"github.com/femitubosun/go-sweepline-availability/internal/location"
)

const apiV1 = "/api/v1"

func (a *app) registerLocationRoutes(m *http.ServeMux) {
	h := location.NewHandler(a.services.location)

	m.Handle("GET "+apiV1+"/default-location", http.HandlerFunc(h.GetLocationInfo))
}

func (a *app) registerBookingRoutes(m *http.ServeMux) {
	h := booking.NewHandler(a.services.booking)

	m.Handle("GET "+apiV1+"/courts/{courtID}/bookings", http.HandlerFunc(h.GetCourtBookings))
	m.Handle("POST "+apiV1+"/courts/{courtID}/bookings", http.HandlerFunc(h.CreateBooking))
	m.Handle("GET "+apiV1+"/bookings/{bookingID}", http.HandlerFunc(h.GetBooking))
	m.Handle("PATCH "+apiV1+"/bookings/{bookingID}/cancel", http.HandlerFunc(h.CancelBooking))
	m.Handle("GET "+apiV1+"/courts/{courtID}/availability", http.HandlerFunc(h.AvailabilitySearch))
}

func (a *app) registerStaticRoute(m *http.ServeMux) {
	m.Handle("GET /", http.FileServer(http.Dir("static")))
}
