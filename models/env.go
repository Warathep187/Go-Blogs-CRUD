package models

type EnvVar struct {
	AppEnv string
	Port   int

	MongoEndpoint string
	MongoUsername string
	MongoPassword string
	MongoDatabase string
}
