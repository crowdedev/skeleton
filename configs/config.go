package configs

func LoadConfigs() {
	_LoadEnv()
	_LoadDatabase()
	_LoadElasticsearch()
}
