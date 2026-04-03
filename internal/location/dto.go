package location

import (
	"fmt"
	"time"
)

type LocationResponse struct {
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	Timezone      string                     `json:"timezone"`
	Hours         map[string][]HoursResponse `json:"hours"`
	Exceptions    []ExceptionResponse        `json:"exceptions"`
	BufferMinutes int                        `json:"bufferMinutes"`
	Courts        []CourtResponse            `json:"courts"`
}

type HoursResponse struct {
	Open  string `json:"open"`
	Close string `json:"close"`
	Range string `json:"range"`  // e.g., "08:00-20:00"
}

type ExceptionResponse struct {
	Start  string `json:"start"`
	End    string `json:"end"`
	Reason string `json:"reason"`
}

type CourtResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

func toLocationResponse(loc Location) LocationResponse {
	return LocationResponse{
		ID:            loc.ID,
		Name:          loc.Name,
		Timezone:      loc.Timezone,
		BufferMinutes: loc.BufferMinutes,
		Hours:         toHoursResponse(loc.Hours),
		Exceptions:    toExceptionsResponse(loc.Exceptions),
		Courts:        toCourtsResponse(loc.Courts),
	}
}

func toHoursResponse(hours OperatingHours) map[string][]HoursResponse {
	result := make(map[string][]HoursResponse)

	for weekday, schedule := range hours.Schedule {
		dayName := weekday.String()
		var dayHours []HoursResponse

		for _, r := range schedule.Ranges {
			openStr := formatDuration(r.Open)
			closeStr := formatDuration(r.Close)
			dayHours = append(dayHours, HoursResponse{
				Open:  openStr,
				Close: closeStr,
				Range: openStr + "-" + closeStr,
			})
		}

		result[dayName] = dayHours
	}

	return result
}

func toCourtsResponse(courts []Court) []CourtResponse {
	result := make([]CourtResponse, len(courts))

	for i, c := range courts {
		result[i] = CourtResponse{
			ID:       c.ID,
			Name:     c.Name,
			Capacity: c.Capacity,
		}
	}

	return result
}

func toExceptionsResponse(exceptions []Exception) []ExceptionResponse {
	result := make([]ExceptionResponse, len(exceptions))

	for i, e := range exceptions {
		result[i] = ExceptionResponse{
			Start:  e.Start.Format(time.RFC3339),
			End:    e.End.Format(time.RFC3339),
			Reason: e.Reason,
		}
	}

	return result
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
