package configs

import "github.com/sirupsen/logrus"

type (
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

	LoggerExtension struct {
		Extensions []logrus.Hook
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
