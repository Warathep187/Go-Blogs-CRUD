package configs

import (
	models "go_blogs/models"

	"github.com/spf13/viper"
)

var Env models.EnvVar

func setDefaultConfig(v *viper.Viper) {
	defaultPort := 8080
	v.SetDefault("PORT", defaultPort)
}

func InitEnv() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	setDefaultConfig(viper.GetViper())

	Env.AppEnv = viper.GetString("APP_ENV")
	Env.Port = viper.GetInt("PORT")

	Env.MongoEndpoint = viper.GetString("MONGO_ENDPOINT")
	Env.MongoUsername = viper.GetString("MONGO_USERNAME")
	Env.MongoPassword = viper.GetString("MONGO_PASSWORD")
	Env.MongoDatabase = viper.GetString("MONGO_DATABASE")
}
