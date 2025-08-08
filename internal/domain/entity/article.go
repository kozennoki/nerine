package entity

import "time"

type Article struct {
	ID          string
	Title       string
	Image       string
	Category    Category
	Description string
	Body        string
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
