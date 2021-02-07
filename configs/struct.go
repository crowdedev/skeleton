package configs

import (
	"time"

	uuid "github.com/google/uuid"
	gorm "gorm.io/gorm"
)

type (
	Base struct {
		Id        uuid.UUID      `gorm:"type:varchar(36);primaryKey"`
		CreatedAt time.Time      `gorm:"type:timestamp(3);notNull;default:CURRENT_TIMESTAMP(3)"`
		UpdatedAt time.Time      `gorm:"type:timestamp(3);notNull;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP"`
		CreatedBy string         `gorm:"type:varchar(36);default:null"`
		UpdatedBy string         `gorm:"type:varchar(36);default:null"`
		DeletedAt gorm.DeletedAt `gorm:"default:null;index"`
		DeletedBy string         `gorm:"type:varchar(36);default:null"`
	}

	User struct {
		Id    string
		Email string
		Role  string
	}

	Env struct {
		Debug              bool
		HtppPort           int
		RpcPort            int
		Version            string
		ApiVersion         string
		ServiceName        string
		DbHost             string
		DbPort             int
		DbUser             string
		DbPassword         string
		DbName             string
		DbDriver           string
		DbAutoMigrate      bool
		ElasticsearchHost  string
		ElasticsearchPort  int
		ElasticsearchIndex string
		MongoDbHost        string
		MongoDbPort        int
		MongoDbName        string
		AmqpHost           string
		AmqpPort           int
		AmqpUser           string
		AmqpPassword       string
		HeaderUserId       string
		HeaderUserEmail    string
		HeaderUserRole     string
		CacheLifetime      int
		User               *User
		TemplateLocation   string
	}

	Template struct {
		ApiVersion            string
		PackageName           string
		Module                string
		ModuleLowercase       string
		ModulePlural          string
		ModulePluralLowercase string
		Columns               []*FieldTemplate
	}

	ModuleTemplate struct {
		Name   string
		Fields []*FieldTemplate
	}

	FieldTemplate struct {
		Name           string
		NameUnderScore string
		Type           string
		Index          int
		IsRequired     bool
	}
)
