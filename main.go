package main

import (
	"fmt"
	"log"

	"github.com/gioCuesta25/employees-manager-backend/api"
	"github.com/gioCuesta25/employees-manager-backend/config"
	"github.com/gioCuesta25/employees-manager-backend/database"
)

func main() {
	env, err := config.LoadEnvironment()

	if err != nil {
		log.Fatal("Error loading env variables")
	}

	db, err := database.NewDbConnection(env)

	if err != nil {
		log.Fatal("Error connecting to database: ", err.Error())
	}

	server := api.NewServer(env.ApiPort, db)

	server.Run()

	fmt.Println(env)
}
