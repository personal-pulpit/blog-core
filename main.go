package main

import (
	"blog/api/server"
	"blog/data/db"
)

func main() {
	db.InitDB()
	server.InitServer()
}