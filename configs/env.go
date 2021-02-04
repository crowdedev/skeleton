package configs

type Env struct {
	Debug              bool
	HtppPort           int
	RpcPort            int
	Version            string
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
}
