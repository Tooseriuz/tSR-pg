package domain

import "time"

type Journey struct {
	ID        int64
	Name      string
	Timestamp time.Time
	Location  string
	Thumbnail *string
}

type JourneyContent struct {
	Name      string
	Timestamp time.Time
	Content   string
}

type CreateJourney struct {
	Name      string
	Location  string
	Thumbnail *string
	Content   string
	Highlight bool
}
