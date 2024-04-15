package main

import (
	"context"
	"fmt"
	"rent-car/api"
	"rent-car/config"
	"rent-car/pkg/logger"
	"rent-car/service"
	"rent-car/storage/postgres"
	"rent-car/storage/redis"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	newRedis := redis.New(cfg)
	store, err := postgres.New(context.Background(), cfg,log,newRedis)
	if err != nil {
		fmt.Println("error while connecting db, err: ", err)
		return
	}
	defer store.CloseDB()

	services := service.New(store,log,newRedis)
	c := api.New(services, log)

	fmt.Println("programm is running on localhost:8080...")
	c.Run(":8080")
}
