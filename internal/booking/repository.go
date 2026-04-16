package booking

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	appRedis "github.com/femitubosun/go-sweepline-availability/internal/redis"
)

type Interval struct {
	Start time.Time
	End   time.Time
}

type Repository interface {
	Reserve(ctx context.Context, input ReserveInput) (Booking, error)
	Get(ctx context.Context, bookingID string) (Booking, error)
	Release(ctx context.Context, bookingID string) error
	GetGaps(ctx context.Context, input GetGapsInput) error
	ListByCourtId(ctx context.Context, courtID string) ([]Booking, error)
}

type repo struct {
	cache *appRedis.Cache
}

func (r *repo) Reserve(ctx context.Context, input ReserveInput) (Booking, error) {
	keys := []string{
		fmt.Sprintf("bookings:{%s}", input.CourtID),
		"booking",
	}
	args := []any{
		input.ID,
		input.Interval.Start.UnixMilli(),
		input.Interval.End.UnixMilli(),
		input.GuestName,
	}

	result, err := reserveScript.Run(ctx, r.cache.Client(), keys, args...).Int()
	if err != nil {
		return Booking{}, fmt.Errorf("redis script failed: %w", err)
	}

	if result == 0 {
		return Booking{}, ErrConflict
	}

	return Booking{
		ID:        input.ID,
		CourtID:   input.CourtID,
		Interval:  input.Interval,
		GuestName: input.GuestName,
	}, nil
}

func (r *repo) Get(ctx context.Context, bookingID string) (Booking, error) {
	keys := []string{
		fmt.Sprintf("booking:%s", bookingID),
	}

	result, err := getBookingScript.Run(ctx, r.cache.Client(), keys).Result()
	if err != nil {
		return Booking{}, fmt.Errorf("redis script failed: %w", err)
	}

	flatArray, ok := result.([]any)
	if !ok {
		return Booking{}, fmt.Errorf("unexpected result type: %T", result)
	}

	if len(flatArray) == 0 {
		return Booking{}, ErrNotFound
	}

	booking := Booking{}

	for i := 0; i < len(flatArray); i += 2 {
		if i+1 >= len(flatArray) {
			break
		}

		key := flatArray[i].(string)
		value := flatArray[i+1].(string)

		switch key {
		case "id":
			booking.ID = value
		case "start":
			startMs, _ := strconv.ParseInt(value, 10, 64)
			booking.Interval.Start = time.UnixMilli(startMs)
		case "end":
			endMs, _ := strconv.ParseInt(value, 10, 64)
			booking.Interval.End = time.UnixMilli(endMs)
		case "guestName":
			booking.GuestName = value
		case "courtId":
			booking.CourtID = value
		}
	}

	return booking, nil
}

func (r *repo) Release(ctx context.Context, bookingID string) error {
	keys := []string{
		fmt.Sprintf("booking:%s", bookingID),
	}
	args := []any{
		bookingID,
	}

	result, err := releaseBookingScript.Run(ctx, r.cache.Client(), keys, args).Int()
	if err != nil {
		return fmt.Errorf("redis script failed: %w", err)
	}

	if result == 0 {
		return ErrNotFound
	}

	return nil

}

func (r *repo) ListByCourtId(ctx context.Context, courtID string) ([]Booking, error) {
	key := fmt.Sprintf("bookings:{%s}", courtID)

	bookingIDs, err := r.cache.Client().ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get booking ids: %w", err)
	}

	if len(bookingIDs) == 0 {
		return make([]Booking, 0), nil
	}
	pipe := r.cache.Client().Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(bookingIDs))

	for i, bookingID := range bookingIDs {
		cmds[i] = pipe.HGetAll(ctx, fmt.Sprintf("booking:%s", bookingID))
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("pipeline exec failed: %w", err)
	}

	bookings := make([]Booking, 0, len(bookingIDs))
	for i, cmd := range cmds {
		result, err := cmd.Result()
		if err != nil {
			slog.Error("HGETALL returned error", "index", i, "bookingID", bookingIDs[i], "error", err)
		}

		if len(result) == 0 {
			slog.Warn("HGETALL returned empty hash", "index", i, "bookingID", bookingIDs[i])
			continue
		}

		booking := Booking{
			ID:      bookingIDs[i],
			CourtID: courtID,
		}

		if startStr, ok := result["start"]; ok {
			startMs, _ := strconv.ParseInt(startStr, 10, 64)
			booking.Interval.Start = time.UnixMilli(startMs)
		}
		if endStr, ok := result["end"]; ok {
			endMs, _ := strconv.ParseInt(endStr, 10, 64)
			booking.Interval.End = time.UnixMilli(endMs)
		}
		if guestName, ok := result["guestName"]; ok {
			booking.GuestName = guestName
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *repo) GetGaps(ctx context.Context, input GetGapsInput) error {
	return nil
}

func NewRepo(cache *appRedis.Cache) Repository {
	return &repo{
		cache: cache,
	}
}

type ReserveInput struct {
	ID        string
	CourtID   string
	Interval  Interval
	GuestName string
}

type GetGapsInput struct {
	CourtID      string
	DurationMins int
	Interval     Interval
}

var ErrConflict = fmt.Errorf("booking conflicts with existing reservation")
var ErrNotFound = fmt.Errorf("booking not found")
