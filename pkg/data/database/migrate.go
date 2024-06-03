package database

import "blog/pkg/logging"

func Migration(models any) {
	err := DB.AutoMigrate(&models)
	if err != nil {
		logging.MyLogger.Fatal(logging.Mysql, logging.Migration, err.Error(), nil)

	}
}
