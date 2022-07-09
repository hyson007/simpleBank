package util

import (
	"time"

	"github.com/spf13/viper"
)

//Config stores all configuration of the app
//The values are read by viper from config file or env variable
type Config struct {
	DBDriver            string        `mapstructure:"DBDRIVER"`
	DBSource            string        `mapstructure:"DBSOURCE"`
	ServerAddress       string        `mapstructure:"SERVERADDRESS"`
	TokenSymKey         string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // can use any other file, such as json, xml
	viper.AutomaticEnv()       // also read from env
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
