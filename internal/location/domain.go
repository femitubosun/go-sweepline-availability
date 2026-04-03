package location

import "time"

type Location struct {
	ID            string
	Name          string
	Timezone      string
	Hours         OperatingHours
	Exceptions    []Exception
	BufferMinutes int
	Courts        []Court
}

type OperatingHours struct {
	Schedule map[time.Weekday]DaySchedule
}

type DaySchedule struct {
	Ranges []TimeRange
}

type TimeRange struct {
	Open  time.Duration // e.g., 9h * time.Hour = 9:00 AM
	Close time.Duration // e.g., 17h * time.Hour = 5:00 PM
}

type Exception struct {
	Start  time.Time // Full datetime (date + time)
	End    time.Time // Full datetime (date + time)
	Reason string
}

func (e Exception) IsClosedOn(t time.Time) bool {
	return !t.Before(e.Start) && t.Before(e.End)
}

func (oh OperatingHours) IsOpenAt(t time.Time) bool {
	weekday := t.Weekday()
	daySchedule, exists := oh.Schedule[weekday]
	if !exists {
		return false
	}

	minutesFromMidnight := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute

	for _, r := range daySchedule.Ranges {
		if minutesFromMidnight >= r.Open && minutesFromMidnight < r.Close {
			return true
		}
	}
	return false
}

type Court struct {
	ID       string
	Name     string
	Capacity int
}

func NewDefaultLocation() Location {
	return Location{
		ID:            "regents-001",
		Name:          "Regent's Park Tennis Centre",
		Timezone:      "Europe/London",
		BufferMinutes: 10,
		Hours: OperatingHours{
			Schedule: map[time.Weekday]DaySchedule{
				time.Monday:    {Ranges: []TimeRange{{8 * time.Hour, 20 * time.Hour}}},
				time.Tuesday:   {Ranges: []TimeRange{{8 * time.Hour, 20 * time.Hour}}},
				time.Wednesday: {Ranges: []TimeRange{{8 * time.Hour, 20 * time.Hour}}},
				time.Thursday:  {Ranges: []TimeRange{{8 * time.Hour, 20 * time.Hour}}},
				time.Friday:    {Ranges: []TimeRange{{8 * time.Hour, 18 * time.Hour}}},
				time.Saturday:  {Ranges: []TimeRange{{9 * time.Hour, 17 * time.Hour}}},
				time.Sunday:    {Ranges: []TimeRange{{9 * time.Hour, 16 * time.Hour}}},
			},
		},
		Exceptions: []Exception{
			{
				Start:  time.Date(2026, 4, 15, 8, 0, 0, 0, time.UTC),
				End:    time.Date(2026, 4, 17, 18, 0, 0, 0, time.UTC),
				Reason: "Court resurfacing maintenance",
			},
			{
				Start:  time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC),
				End:    time.Date(2026, 5, 4, 23, 59, 59, 0, time.UTC),
				Reason: "Early May Bank Holiday",
			},
			{
				Start:  time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC),
				End:    time.Date(2026, 6, 14, 18, 0, 0, 0, time.UTC),
				Reason: "Summer Championship",
			},
		},
		Courts: []Court{
			{ID: "court-001", Name: "Centre Court", Capacity: 4},
			{ID: "court-002", Name: "Court 2", Capacity: 4},
			{ID: "court-003", Name: "Court 3", Capacity: 4},
			{ID: "court-004", Name: "Court 4", Capacity: 2},
			{ID: "court-005", Name: "Court 5", Capacity: 2},
		},
	}
}
