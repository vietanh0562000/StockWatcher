package main

import (
	"log"
	"stockwatcher/config"
	"stockwatcher/database"
)

func main() {
	config := config.LoadConfig()

	db, err := database.Connect(config)
	if err != nil {
		log.Fatal("Failed to connect database: ", err.Error())
		return
	}

	database.GetStorageInstance(db)

}
