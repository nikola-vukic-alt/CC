package main

import (
	"context"
	"library-app/local/controller"
	"library-app/local/repository"
	"library-app/local/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	borrows_database := client.Database("borrows")
	borrowRepo := repository.NewBorrowRepository(borrows_database)
	borrowService := service.NewBorrowService(borrowRepo)
	borrowController := controller.NewBorrowController(borrowService)

	router := http.NewServeMux()
	router.HandleFunc("/borrow", borrowController.BorrowBook)

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		log.Println("Local library server listening on :8081")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan

	log.Println("Shutting down the local library server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	log.Println("Server gracefully stopped")
}
