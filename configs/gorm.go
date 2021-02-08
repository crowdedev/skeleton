package configs

import (
	"time"

	uuid "github.com/google/uuid"
	gorm "gorm.io/gorm"
)

type Base struct {
	Id        uuid.UUID      `gorm:"type:varchar(36);primaryKey"`
	CreatedAt time.Time      `gorm:"type:timestamp(3);notNull;default:CURRENT_TIMESTAMP(3)"`
	UpdatedAt time.Time      `gorm:"type:timestamp(3);notNull;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP"`
	CreatedBy string         `gorm:"type:varchar(36);default:null"`
	UpdatedBy string         `gorm:"type:varchar(36);default:null"`
	DeletedAt gorm.DeletedAt `gorm:"default:null;index"`
	DeletedBy string         `gorm:"type:varchar(36);default:null"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Id == uuid.Nil {
		b.Id = uuid.New()
	}

	return nil
}
