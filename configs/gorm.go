package configs

import (
	"database/sql"

	uuid "github.com/google/uuid"
	gorm "gorm.io/gorm"
)

type Base struct {
	ID        string
	Counter   uint64 `gorm:"primaryKey;autoIncrement:true"`
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	CreatedBy sql.NullString
	UpdatedBy sql.NullString
	DeletedAt gorm.DeletedAt
	DeletedBy sql.NullString
}

// Get/Set

func (b *Base) SetCreatedBy(user *User) {
	m.CreatedBy = sql.NullString{String: user.Id, Valid: true}
}

func (b *Base) SetUpdatedBy(user *User) {
	m.UpdatedBy = sql.NullString{String: user.Id, Valid: true}
}

func (b *Base) SetDeletedBy(user *User) {
	m.DeletedBy = sql.NullString{String: user.Id, Valid: true}
}

func (b *Base) SetCreatedAt(time time.Time) {
	m.CreatedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *Base) SetUpdatedAt(time time.Time) {
	m.UpdatedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *Base) SetDeletedAt(time time.Time) {
	m.DeletedAt = gorm.DeletedAt{Time: time, Valid: true}
}

// Hooks

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()

	return nil
}
