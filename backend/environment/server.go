package environment

import "github.com/kelseyhightower/envconfig"

type EnvironmentVariables struct {
	UsersDbUser     string `envconfig:"USERS_DB_USER"`
	UsersDbPassword string `envconfig:"USERS_DB_PASSWORD"`
	UsersDbHost     string `envconfig:"USERS_DB_HOST"`
	ItemsDbPassword string `envconfig:"ITEMS_DB_PASSWORD"`
	ItemsDbUser     string `envconfig:"ITEMS_DB_USER"`
	ItemsDbHost     string `envconfig:"ITEMS_DB_HOST"`
	KafkaHost       string `envconfig:"KAFKA_HOST" default:""`
	KafkaUser       string `envconfig:"KAFKA_USER" default:""`
	KafkaPassword   string `envconfig:"KAFKA_PASSWORD" default:""`
	RabbitQueue     string `envconfig:"RABBIT_QUEUE" default:"my-queue"`
	RabbitUrl       string `envconfig:"RABBIT_URL" default:"amqp://user:password@localhost:5672/"`
	RabbitUser      string `envconfig:"KAFKA_PASSWORD" default:"user"`
	RabbitPassword  string `envconfig:"RABBIT_USER" default:""`
	ElasticHost     string `envconfig:"ELASTIC_HOST" default:"http://localhost:9200"`
	ElasticUsername string `envconfig:"ELASTIC_USER_NAME" default:"none"`
	ElasticPassword string `envconfig:"ELASTIC_PASSWORD" default:"none"`
	AwsRegion       string `envconfig:"AWS_REGION" default:"none"`
	RedisHost       string `envconfig:"REDIS_HOST" default:"none"`
	RedisUser       string `envconfig:"REDIS_USER" default:"none"`
	RedisPassword   string `envconfig:"REDIS_PASSWORD" default:"none"`
}

var Variables EnvironmentVariables

func Load() error {
	err := envconfig.Process("", &Variables)

	if err != nil {
		return err
	}

	return nil
}
