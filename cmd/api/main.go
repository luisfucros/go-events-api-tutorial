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
	config		  configs.Config
	store         store.Storage
	cacheStorage  cache.Storage
	logger       *zap.SugaredLogger
}

func main() {

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// DB
	cfg := mysqlDriver.Config{
		User:                 configs.Envs.DB.User,
		Passwd:               configs.Envs.DB.Password,
		Addr:                 configs.Envs.DB.Address,
		DBName:               configs.Envs.DB.Name,
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
	if configs.Envs.Redis.Enabled {
		rdb = cache.NewRedisClient(
			configs.Envs.Redis.Addr,
			configs.Envs.Redis.Password,
			configs.Envs.Redis.DB,
		)
		logger.Info("redis cache connection established")
		defer rdb.Close()
	}
	
	// Storage
	storage := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	app := &application{
		config:       configs.Envs,
		store:        storage,
		cacheStorage: cacheStorage,
		logger:       logger,
	}

	if err := app.serve(); err != nil {
		logger.Fatal(err)
	}
}