package qentity

import "time"

type TableUsage struct {
	ID        string     `json:"id"`
	CompanyID string     `json:"company_id" binding:"required"`
	Key       string     `json:"key" binding:"required"`
	Value     int        `json:"value"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by" binding:"required"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy string     `json:"updated_by" binding:"required"`
	DeletedAt NullTime   `json:"deleted_at"`
	DeletedBy NullString `json:"deleted_by"`
}
