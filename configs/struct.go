package configs

import (
	"time"

	"gorm.io/gorm"
)

type (
	Base struct {
		Id        string         `gorm:"type:varchar(36);primary_key"`
		CreatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time      `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
		CreatedBy string         `gorm:"type:varchar(36);default:null"`
		UpdatedBy string         `gorm:"type:varchar(36);default:null"`
		DeletedAt gorm.DeletedAt `gorm:"default:null;index"`
		DeletedBy string         `gorm:"type:varchar(36);default:null"`
	}

	User struct {
		ID    string
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
		ProtobufType   string
		GolangType     string
		Index          int
		IsRequired     bool
	}
)
