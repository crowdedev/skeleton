package configs

import (
	"database/sql"

	uuid "github.com/google/uuid"
	gorm "gorm.io/gorm"
)

type Base struct {
	ID        string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	CreatedBy sql.NullString
	UpdatedBy sql.NullString
	DeletedAt gorm.DeletedAt
	DeletedBy sql.NullString
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()

	return nil
}
