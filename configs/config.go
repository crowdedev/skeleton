package configs

func LoadConfigs() {
	loadEnv()
	loadDatabase()
	loadElasticsearch()
}
