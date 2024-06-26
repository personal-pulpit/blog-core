package database

import "blog/pkg/logging"

func Migration(model any) {
	err := postgresInstance.AutoMigrate(&model)
	if err != nil {
		logging.MyLogger.Fatal(logging.Postgres, logging.Migration, err.Error(), nil)

	}
}
