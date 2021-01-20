package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	env struct {
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
		User               *User
	}
)

var Env env

func loadEnv() {
	godotenv.Load()

	Env.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
	Env.HtppPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	Env.RpcPort, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))

	Env.DbHost = os.Getenv("DB_HOST")
	Env.DbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	Env.DbUser = os.Getenv("DB_USER")
	Env.DbPassword = os.Getenv("DB_PASSWORD")
	Env.DbName = os.Getenv("DB_NAME")
	Env.DbAutoMigrate, _ = strconv.ParseBool(os.Getenv("DB_AUTO_CREATE"))

	Env.ElasticsearchHost = os.Getenv("ELASTICSEARCH_HOST")
	Env.ElasticsearchPort, _ = strconv.Atoi(os.Getenv("ELASTICSEARCH_PORT"))
	Env.ElasticsearchIndex = Env.DbName

	Env.MongoDbHost = os.Getenv("MONGODB_HOST")
	Env.MongoDbPort, _ = strconv.Atoi(os.Getenv("MONGODB_PORT"))
	Env.MongoDbName = os.Getenv("MONGODB_NAME")

	Env.AmqpHost = os.Getenv("AMQP_HOST")
	Env.AmqpPort, _ = strconv.Atoi(os.Getenv("AMQP_PORT"))
	Env.AmqpUser = os.Getenv("AMQP_USER")
	Env.AmqpPassword = os.Getenv("AMQP_PASSWORD")

	Env.HeaderUserId = os.Getenv("HEADER_USER_ID")
	Env.HeaderUserEmail = os.Getenv("HEADER_USER_EMAIL")
	Env.HeaderUserRole = os.Getenv("HEADER_USER_ROLE")

	Env.User = &User{}
}
