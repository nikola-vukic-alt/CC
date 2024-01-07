package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"library-app/central/controller"
	"library-app/central/repository"
	"library-app/central/service"

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

	database := client.Database("members")
	memberRepo := repository.NewMemberRepository(database)
	memberService := service.NewMemberService(memberRepo)
	memberController := controller.NewMemberController(memberService)

	router := http.NewServeMux()
	router.HandleFunc("/register", memberController.Register)
	router.HandleFunc("/get", memberController.GetMemberBySSN)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Central library server listening on :8080")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan

	log.Println("Shutting down the central library server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	log.Println("Server gracefully stopped")
}
