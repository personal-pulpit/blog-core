package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)



type (
	Config struct {
		Redis  MyRedisConfig
		Mysql  MysqlConfig
		Jwt    JwtConfig
		Server ServerConfig
		Logger LoggerConfig
	}
	MysqlConfig struct {
		Host      string
		Username  string
		Password  string
		Port      string
		DBname    string
		ParseTime bool
	}
	JwtConfig struct {
		Secret string
	}
	ServerConfig struct {
		Port string
	}
	MyRedisConfig struct {
		Host     string
		Username string
		Password string
		Port     string
		DBname   string
		Protocol string
	}
	LoggerConfig struct {
		LogFilePath string
		LoggerName  string
		Level       string
		Encoding    string
	}
)
var (
	cfg *Config
	mu = &sync.Mutex{}
)
func ReadConfigs()*Config{
	mu.Lock()
	defer mu.Unlock()
	if cfg == nil{
		newConfig := &Config{}
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("../config/")
		viper.AutomaticEnv()
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %v", err))
		}
		if err := viper.Unmarshal(&cfg); err != nil {
			panic(fmt.Errorf("error unmarshaling config: %s", err))
			
		}
		cfg = newConfig
	}
	return cfg
	
}
