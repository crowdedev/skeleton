package configs

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	gorm "gorm.io/gorm"
)

type Base struct {
	ID        string `gorm:"primaryKey;autoIncrement:false"`
	Counter   uint64 `gorm:"primaryKey;autoIncrement:true"`
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	CreatedBy sql.NullString
	UpdatedBy sql.NullString
	DeletedAt gorm.DeletedAt
	DeletedBy sql.NullString
}

func (b *Base) SetCreatedBy(user *User) {
	b.CreatedBy = sql.NullString{String: user.Id, Valid: true}
}

func (b *Base) SetUpdatedBy(user *User) {
	b.UpdatedBy = sql.NullString{String: user.Id, Valid: true}
}

func (b *Base) SetDeletedBy(user *User) {
	b.DeletedBy = sql.NullString{String: user.Id, Valid: true}
}

func (b *Base) SetCreatedAt(time time.Time) {
	b.CreatedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *Base) SetUpdatedAt(time time.Time) {
	b.UpdatedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *Base) SetDeletedAt(time time.Time) {
	b.DeletedAt = gorm.DeletedAt{Time: time, Valid: true}
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()

	return nil
}
