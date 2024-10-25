package qentity

import "time"

type TableLimit struct {
	ID        string     `json:"id"`
	CompanyID string     `json:"company_id" binding:"required"`
	Key       string     `json:"key" binding:"required"`
	Value     int        `json:"value"`
	Type      int        `json:"type" binding:"required"`
	Name      string     `json:"name"`
	Metric    string     `json:"metric"`
	Desc      string     `json:"desc"`
	ExpiredAt time.Time  `json:"expired_at"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by" binding:"required"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy string     `json:"updated_by" binding:"required"`
	DeletedAt NullTime   `json:"deleted_at"`
	DeletedBy NullString `json:"deleted_by"`
}
