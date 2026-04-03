package main

import (
	"net/http"

	"github.com/femitubosun/go-sweepline-availability/internal/availability"
	"github.com/femitubosun/go-sweepline-availability/internal/location"
)

const apiV1 = "/api/v1"

func (a *app) registerLocationRoutes(m *http.ServeMux) {
	h := location.NewHandler(a.services.location)

	m.Handle("GET "+apiV1+"/default-location", http.HandlerFunc(h.GetLocationInfo))
}

func (a *app) registerAvailabilityRoutes(m *http.ServeMux) {
	h := availability.NewHandler(a.services.availability)

	m.Handle("GET "+apiV1+"/courts/{id}/availability", http.HandlerFunc(h.GetAvailability))
}

func (a *app) registerStaticRoute(m *http.ServeMux) {
	m.Handle("GET /", http.FileServer(http.Dir("static")))
}
