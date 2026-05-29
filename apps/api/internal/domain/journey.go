package domain

import "time"

type Journey struct {
	Name      string
	Timestamp time.Time
	Location  string
	Thumbnail string
}
