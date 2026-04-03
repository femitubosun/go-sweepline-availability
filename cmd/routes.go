package main

import (
	"net/http"

	"github.com/femitubosun/go-sweepline-availability/internal/location"
)

func (a *app) registerLocationRoutes(m *http.ServeMux) {
	h := location.NewHandler(a.services.location)

	m.Handle("GET /location", http.HandlerFunc(h.GetLocationInfo))
}

func (a *app) registerStaticRoute(m *http.ServeMux) {
	m.Handle("GET /", http.FileServer(http.Dir("static")))
}
