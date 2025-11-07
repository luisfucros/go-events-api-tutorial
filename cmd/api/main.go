package main

import (
	"log"
	"github.com/luisfucros/go-events-api-tutorial/internal/database"
	"github.com/luisfucros/go-events-api-tutorial/internal/configs"
	_ "github.com/luisfucros/go-events-api-tutorial/docs"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

// @title Events Rest API
// @version 0.1.0
// @description	API for event and users
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

type application struct {
	port       int64
	JWTSecret  string
	models     database.Models

}

func main() {
	cfg := mysqlDriver.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := database.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	models := database.NewModels(db)

	app := &application{
		port: configs.Envs.Port,
		JWTSecret: configs.Envs.JWTSecret,
		models: models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}