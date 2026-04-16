package booking

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"
)

var ErrHistoricalBooking = errors.New("booking cannot be in the past")

type Service interface {
	CreateBooking(ctx context.Context, input CreateBookingInput) (Booking, error)
	ListCourtBookings(ctx context.Context, courtID string) ([]Booking, error)
	GetBooking(ctx context.Context, bookingID string) (Booking, error)
	GetGaps(ctx context.Context, input GetGapsInput) ([]Interval, error)
	CancelBooking(ctx context.Context, bookingId string) error
}

type svc struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) CreateBooking(ctx context.Context, input CreateBookingInput) (Booking, error) {
	if input.StartTime.Before(time.Now()) {
		return Booking{}, ErrHistoricalBooking
	}

	endTime := input.StartTime.Add(time.Duration(input.DurationMinutes) * time.Minute)

	booking, err := s.repo.Reserve(ctx, ReserveInput{
		ID:        generateBookingId(),
		CourtID:   input.CourtID,
		Interval:  Interval{Start: input.StartTime, End: endTime},
		GuestName: input.GuestName,
	})
	if err != nil {
		if err == ErrConflict {
			return Booking{}, fmt.Errorf("slot unavailable: %w", err)
		}
		return Booking{}, fmt.Errorf("failed to create booking: %w", err)
	}

	return booking, nil
}

func (s *svc) ListCourtBookings(ctx context.Context, courtID string) ([]Booking, error) {
	bookings, err := s.repo.ListByCourtId(ctx, courtID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}
	return bookings, nil
}

func (s *svc) GetBooking(ctx context.Context, bookingID string) (Booking, error) {
	booking, err := s.repo.Get(ctx, bookingID)

	if err != nil {
		return Booking{}, fmt.Errorf("failed to get bookings: %w", err)
	}
	return booking, nil

}

func (s *svc) CancelBooking(ctx context.Context, bookingID string) error {
	err := s.repo.Release(ctx, bookingID)

	if err != nil {
		return fmt.Errorf("booking does not exist: %w", err)
	}
	return nil
}

func (s *svc) GetGaps(ctx context.Context, input GetGapsInput) ([]Interval, error) {
	bookings, err := s.repo.ListByCourtId(ctx, input.CourtID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	var overlapping []Booking
	for _, b := range bookings {
		if b.Interval.End.After(input.Interval.Start) && b.Interval.Start.Before(input.Interval.End) {
			overlapping = append(overlapping, b)
		}
	}

	sort.Slice(overlapping, func(i, j int) bool {
		return overlapping[i].Interval.Start.Before(overlapping[j].Interval.Start)
	})

	var gaps []Interval
	duration := time.Duration(input.DurationMins) * time.Minute

	if len(overlapping) == 0 {

		if input.Interval.End.Sub(input.Interval.Start) >= duration {
			gaps = append(gaps, input.Interval)
		}
		return gaps, nil
	}

	firstBooking := overlapping[0]
	if firstBooking.Interval.Start.After(input.Interval.Start) {
		gap := Interval{Start: input.Interval.Start, End: firstBooking.Interval.Start}
		if gap.End.Sub(gap.Start) >= duration {
			gaps = append(gaps, gap)
		}
	}

	for i := 0; i < len(overlapping)-1; i++ {
		current := overlapping[i]
		next := overlapping[i+1]

		gap := Interval{Start: current.Interval.End, End: next.Interval.Start}
		if gap.End.Sub(gap.Start) >= duration {
			gaps = append(gaps, gap)
		}
	}

	lastBooking := overlapping[len(overlapping)-1]
	if lastBooking.Interval.End.Before(input.Interval.End) {
		gap := Interval{Start: lastBooking.Interval.End, End: input.Interval.End}
		if gap.End.Sub(gap.Start) >= duration {
			gaps = append(gaps, gap)
		}
	}

	return gaps, nil
}

type CreateBookingInput struct {
	CourtID         string    `validate:"required"`
	StartTime       time.Time `validate:"required"`
	DurationMinutes int       `validate:"required,min=1,max=1440"`
	GuestName       string    `validate:"required"`
}
