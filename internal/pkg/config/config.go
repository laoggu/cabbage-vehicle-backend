package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Cfg struct {
	Env          string   `default:"dev"`
	HttpPort     string   `default:"8080"`
	MysqlDSN     string   `required:"true" split_words:"true"`
	RedisAddr    string   `required:"true" split_words:"true"`
	KafkaBrokers []string `required:"true" split_words:"true"`
	SignSecret   string   `required:"true" split_words:"true"`
	OssEndpoint  string   `required:"true" split_words:"true"`
	OssAK        string   `required:"true" split_words:"true"`
	OssSK        string   `required:"true" split_words:"true"`
}

var C Cfg

func Load() {
	if err := envconfig.Process("cabbage", &C); err != nil {
		log.Fatalf("config load err:%v", err)
	}
}
