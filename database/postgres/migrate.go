package database

func Migration(model any) {
	err := postgresInstance.AutoMigrate(&model)
	if err != nil {
		//TODO:manage panics
		panic(err)
	}
}
