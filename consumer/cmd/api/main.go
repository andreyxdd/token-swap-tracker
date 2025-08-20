package main

import (
	"consumer/internal/rest"
	"consumer/internal/services"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPw := os.Getenv("REDIS_PASSWORD")

	redisCfg := services.RedisConfig{Addr: redisAddr, Password: redisPw}
	repo := services.NewRedisStatsRepo(redisCfg)
	service := services.NewStatsService(repo)

	restApi := rest.New(port, service)
	err := restApi.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
