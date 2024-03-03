package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SharooqSalaudeen/readinsight-backend-golang/config"
	mongoconnection "github.com/SharooqSalaudeen/readinsight-backend-golang/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	// Initialize MongoConnection
	// Create a new instance of MongoConnection
	mongoURL := config.NewConfig().MongoDBURI
	mongoConnection := mongoconnection.NewMongoConnection(mongoURL)

	// Connect to MongoDB
	if mongoURL == "" {
		log.Println("MONGODB_URL not specified in environment")
		os.Exit(1)
	} else {
		mongoConnection.Connect(func() {
			// MongoDB connected successfully
			log.Println("Connected to MongoDB")
		})
	}

	// Close the MongoDB connection on SIGINT signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Gracefully shutting down")
	mongoConnection.Close()
}
