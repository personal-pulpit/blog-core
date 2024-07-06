package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type (
	Config struct {
		Jwt      Jwt      `koanf:"jwt"`
		Server   Server   `koanf:"server"`
		Postgres Postgres `koanf:"postgres"`
		Redis    Redis    `koanf:"redis"`
		Logger   Logger   `koanf:"logger"`
		Email    Email    `koanf:"email"`
	}

	Server struct {
		Port int `koanf:"port"`
	}

	Postgres struct {
		Host     string `koanf:"host"`
		Password string `koanf:"password"`
		Username string `koanf:"username"`
		DBName   string `koanf:"db_name"`
		Port     int    `koanf:"port"`
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
		Username string `koanf:"username"`
		Password string `koanf:"password"`
		Protocol string `koanf:"protocol"`
	}
	Jwt struct {
		Secret string `koanf:"secret"`
	}
	Email struct {
		SenderEmail string `koanf:"sender_email"`
		Password    string `koanf:"password"`
		Host        string `koanf:"host"`
		Port        string `koanf:"port"`
	}
)

var (
	configIns *Config
	mu        = new(sync.Mutex)
)

type Env = string

const (
	Development Env = "development"
	Production  Env = "production"
)

func ConfigsDirPath() string {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error in generating env dir")
	}

	return filepath.Dir(f)
}
func GetConfigInstance() *Config {
	mu.Lock()
	defer mu.Unlock()
	if configIns == nil {
		filename := getConfigFile(GetEnv())

		path := ConfigsDirPath()
		
		k := koanf.New(path)
		
		if err := k.Load(file.Provider(path+"/"+filename), yaml.Parser()); err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		
		var config = &Config{}
		
		if err := k.Unmarshal("", config); err != nil {
			log.Fatalf("error unmarshaling config: %v", err)
		}
		
		configIns = config
	}
	return configIns
}
func GetEnv() Env {
	err := godotenv.Load(ConfigsDirPath()+"/"+".env")
	if err != nil {
		panic(err)
	}

	env := strings.ToLower(os.Getenv("ENV"))

	if env == Development || env == "" {
		return Development
	} else if env == Production {
		return Production
	} else {
		panic("invalid env:" + env)
	}
}
func getConfigFile(env Env) string {
	if env == Development {
		return "config-development.yaml"
	} else if env == Production {
		return "config-production.yaml"
	} else {
		panic("invalid env:" + env)
	}
}
