package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"loadbalancer/db"
	"loadbalancer/routes"
	"log"
	"time"
)

func main() {
	db.IntiRedisDB()

	server := gin.Default()
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(time.Second*10),
		gocron.NewTask(routes.ProcessKeysAndSendRequests),
	)

	if err != nil {
		log.Fatalf("Failed to schedule job: %v", err)
	}

	scheduler.Start()

	routes.Routes(server)

	log.Println("Server running on :1010")
	server.Run(":1010")

}
