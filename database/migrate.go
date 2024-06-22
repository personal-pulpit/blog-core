package database

import "blog/pkg/logging"

func Migration(model any) {
	err := DB.AutoMigrate(&model)
	if err != nil {
		logging.MyLogger.Fatal(logging.Mysql, logging.Migration, err.Error(), nil)

	}
}
