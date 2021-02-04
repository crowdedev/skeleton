package models

import (
	configs "github.com/crowdeco/skeleton/configs"
)

type Child struct {
	configs.Base
	ParentId string `gorm:"column:parentId;type:varchar(36);notNull"`
	Name     string `gorm:"column:name;type:varchar(255);notNull"`
	Parent   Parent `gorm:"constraint:OnDelete:CASCADE"`
}

func (Child) TableName() string {
	return "children"
}

func (j Child) SetCreatedBy(user *configs.User) {
}

func (j Child) SetUpdatedBy(user *configs.User) {
}

func (j Child) SetDeletedBy(user *configs.User) {
}

func (Child) IsSoftDelete() bool {
	return false
}
