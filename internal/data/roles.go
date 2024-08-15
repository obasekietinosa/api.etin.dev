package data

import (
	"time"
)

type Role struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Company     string    `json:"company"`
	CompanyIcon string    `json:"companyIcon"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Skills      []string  `json:"skills"`
}
