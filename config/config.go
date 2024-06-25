package config

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type (
	Config struct {
		Jwt    Jwt    `koanf:"jwt"`
		Server Server `koanf:"server"`
		Mysql  Mysql  `koanf:"mysql"`
		Redis  Redis  `koanf:"redis"`
		Logger Logger `koanf:"logger"`
	}

	Server struct {
		Host string `koanf:"host"`
		Port int    `koanf:"port"`
	}

	Mysql struct {
		Host      string `koanf:"host"`
		Password  string `koanf:"password"`
		Username  string `koanf:"username"`
		DBName    string `koanf:"db_name"`
		Port      int    `koanf:"port"`
		ParseTime bool   `koanf:"prase_time"`
	}
	Logger struct {
		LogFilePath string `koanf:"log_file_path"`
		LoggerName  string `koanf:"logger_name"`
		Level       string `koanf:"level"`
		Encoding    string `koanf:"encoding"`
	}
	Redis struct {
		Host     string `koanf:"host"`
		DB       int    `koanf:"db"`
		Port     int    `koanf:"port"`
		Username string `koanf:"username "`
		Password string `koanf:"password"`
		Protocol string `koanf:"protocol"`
	}
	Jwt struct {
		Secret string `koanf:"secret"`
	}
)

var (
	configIns *Config
	mu        = new(sync.Mutex)
	env       Env
)

type Env string

const (
	Development Env = "development"
	Production  Env = "production"
	Test        Env = "test"
)

func GetConfigInstance() *Config {
	mu.Lock()
	defer mu.Unlock()
	if configIns == nil {
		k := koanf.New("../config")
		if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		var config Config
		if err := k.Unmarshal("", &config); err != nil {
			log.Fatalf("error unmarshaling config: %v", err)
		}
		env = getEnv()
		configIns = &config
	}
	return configIns
}

func getEnv() Env {
	env := strings.ToLower(os.Getenv("ENV"))
	if env == string(Development) || env == "" {
		return Development
	} else if env == string(Production) {
		return Production

	} else if env == string(Test) {
		return Test
	} else {
		panic("invalid env:" + env)
	}
}
