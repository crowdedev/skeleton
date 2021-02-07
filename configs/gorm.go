package configs

import (
	uuid "github.com/google/uuid"
	gorm "gorm.io/gorm"
)

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Id == uuid.Nil {
		b.Id = uuid.New()
	}

	b.CreatedBy = "Env.User.Id"
	b.UpdatedBy = "Env.User.Id"

	return nil
}

func (b *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedBy = "Env.User.Id"

	return nil
}
