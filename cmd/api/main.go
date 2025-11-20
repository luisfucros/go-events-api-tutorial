package main

import (
	"go.uber.org/zap"

	"github.com/luisfucros/go-events-api-tutorial/internal/database"
	"github.com/luisfucros/go-events-api-tutorial/internal/store"
	"github.com/luisfucros/go-events-api-tutorial/internal/configs"
	"github.com/luisfucros/go-events-api-tutorial/internal/store/cache"
	_ "github.com/luisfucros/go-events-api-tutorial/docs"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis/v8"
)

// @title Events Rest API
// @version 0.1.0
// @description	API for event and users
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

type application struct {
	port          int64
	JWTSecret     string
	store         store.Storage
	cacheStorage  cache.Storage
}

func main() {

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// DB
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
		logger.Fatal(err)
	}
	defer db.Close()

	// Redis
	var rdb *redis.Client
	if configs.Envs.REDISEnabled {
		rdb = cache.NewRedisClient(configs.Envs.REDISAddr, configs.Envs.REDISPW, configs.Envs.REDISDB)
		logger.Info("redis cache connection established")

		defer rdb.Close()
	}
	
	// Storage
	storage := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	app := &application{
		port: configs.Envs.Port,
		JWTSecret: configs.Envs.JWTSecret,
		store: storage,
		cacheStorage: cacheStorage,
	}

	if err := app.serve(); err != nil {
		logger.Fatal(err)
	}
}