package configs

import (
	"net/http"
	"time"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type (
	Base struct {
		Id        string         `gorm:"primary_key;column:id;type:varchar(36)"`
		CreatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
		CreatedBy int32          `gorm:"default:null"`
		UpdatedBy int32          `gorm:"default:null"`
		DeletedAt gorm.DeletedAt `gorm:"default:null;index"`
		DeletedBy int32          `gorm:"default:null"`
	}

	Model interface {
		TableName() string
		SetCreatedBy(user *User)
		SetUpdatedBy(user *User)
		SetDeletedBy(user *User)
		IsSoftDelete() bool
	}

	Service interface {
		Name() string
		Create(value interface{}, id string) error
		Update(value interface{}, id string) error
		Bind(value interface{}, id string) error
		Delete(value interface{}, id string) error
	}

	Module interface {
		Consume()
	}

	Server interface {
		RegisterGRpc(server *grpc.Server)
		RegisterAutoMigrate()
		RegisterQueueConsumer()
	}

	Middleware interface {
		Attach(request *http.Request, response http.ResponseWriter) bool
	}

	Application interface {
		Run()
	}
)
