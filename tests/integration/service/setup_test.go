package service

import (
	"blog/config"
	database "blog/database/postgres"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var db *gorm.DB

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
			"POSTGRES_USER":      postgresConfig.Username,
			"POSTGRES_PASSWORD": postgresConfig.Password,
			"POSTGRES_DB":postgresConfig.DBName,
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

	time.Sleep(5*time.Second)
	db, err := database.GetPostgresqlDB(postgresConfig)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TestMain(m *testing.M) {
	db = getDBInstance()

	m.Run()
}

func genEndPort(endpoint string)int{
	endPort, err := strconv.Atoi(strings.Split(endpoint, ":")[1])
	if err != nil {
		log.Fatal(err)
	}

	return endPort
}
