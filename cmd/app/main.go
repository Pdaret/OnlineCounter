package main

import (
	"log"
	"x-ui-monitor/internal/adapter/api"
	logger "x-ui-monitor/internal/adapter/log"
	"x-ui-monitor/internal/repository"
	"x-ui-monitor/internal/usecase"
	"x-ui-monitor/pkg/bolt"
)

const (
	dbFile        = "xui-monitor.db"
	ttlSeconds    = 120 // 2 minutes
	bucketName    = "active_ips"
	accessLogPath = "/path/to/access.log" // Change this
)

func main() {
	db, err := bolt.NewBoltDB(dbFile, bucketName)
	if err != nil {
		log.Fatalf("Failed to initialize BoltDB: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db, bucketName)
	userUsecase := usecase.NewUserUsecase(userRepo)

	go logger.TailLogFile(accessLogPath, userUsecase)

	api.StartServer(userUsecase)
}
