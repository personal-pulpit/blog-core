package database

import "blog/pkg/logging"

func Migration(model any) {
	err := mysqlInstance.AutoMigrate(&model)
	if err != nil {
		logging.MyLogger.Fatal(logging.Mysql, logging.Migration, err.Error(), nil)

	}
}
