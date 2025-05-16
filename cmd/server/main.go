package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"weatherApi/internal/api"
	"weatherApi/internal/db"
	"weatherApi/pkg/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db.ConnectDefaultDB()
	dbInstance := db.DB

	api.SetDB(dbInstance)
	scheduler.SetDB(dbInstance)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	go scheduler.StartWeatherScheduler()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Shutting down gracefully...")
		cancel()
		time.Sleep(2 * time.Second) // дати час scheduler-у завершити роботу
		os.Exit(0)
	}()

	r := gin.Default()
	api.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
