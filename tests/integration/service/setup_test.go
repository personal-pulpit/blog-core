package service

import (
	"blog/config"
	database "blog/database/postgres"
	redisDB "blog/database/redis"

	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var db *gorm.DB
var redisCLI *redis.Client

func getDBInstance() *gorm.DB {
	ctx := context.TODO()

	err := os.Setenv("ENV", "development")
	if err != nil {
		log.Fatalf("Could not set the environment variable to test: %s", err)
	}

	postgresConfig := &config.GetConfigInstance().Postgres

	port := fmt.Sprintf("%d/tcp", postgresConfig.Port)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"POSTGRES_USER":     postgresConfig.Username,
			"POSTGRES_PASSWORD": postgresConfig.Password,
			"POSTGRES_DB":       postgresConfig.DBName,
		},
		WaitingFor: wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(3 * time.Minute),
	}

	postgresContainer, connectErr := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	if connectErr != nil {
		log.Fatal("Failed to start postgres:", connectErr)
	}

	endpoint, err := postgresContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	postgresConfig.Port = genEndPort(endpoint)

	time.Sleep(5 * time.Second)
	db, err := database.GetPostgresqlDB(postgresConfig)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getRedisInstance() *redis.Client {
	ctx := context.TODO()

	err := os.Setenv("ENV", "development")
	if err != nil {
		log.Fatalf("Could not set the environment variable to test: %s", err)
	}

	redisConfig := &config.GetConfigInstance().Redis

	port := fmt.Sprintf("%d/tcp", redisConfig.Port)

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{port},
		Env:          map[string]string{},
		WaitingFor:   wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(3 * time.Minute),
	}

	redisContainer, connectErr := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	if connectErr != nil {
		log.Fatal("Failed to start redis:", connectErr)
	}

	endpoint, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	redisConfig.Port = genEndPort(endpoint)

	time.Sleep(5 * time.Second)
	redis, err := redisDB.GetRedisDB(redisConfig)
	if err != nil {
		log.Fatal(err)
	}

	return redis
}

func TestMain(m *testing.M) {
	db = getDBInstance()
	redisCLI = getRedisInstance()
	m.Run()
}

func genEndPort(endpoint string) int {
	endPort, err := strconv.Atoi(strings.Split(endpoint, ":")[1])
	if err != nil {
		log.Fatal(err)
	}

	return endPort
}
