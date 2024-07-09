package database

func Migration(model any) {
	err := postgresInstance.AutoMigrate(&model)
	if err != nil {
		panic(err)
	}
}
