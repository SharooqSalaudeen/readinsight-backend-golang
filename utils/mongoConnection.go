package mongoconnection

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OnConnectedCallback is a callback for establishing or re-establishing mongo connection
type OnConnectedCallback func()

// MongoConnection is a MongoDB connection wrapper to handle connection issues
type MongoConnection struct {
	mongoURL            string
	onConnectedCallback OnConnectedCallback
	isConnectedBefore   bool
	client              *mongo.Client
}

// NewMongoConnection creates a new MongoConnection instance
func NewMongoConnection(mongoURL string) *MongoConnection {
	return &MongoConnection{
		mongoURL: mongoURL,
	}
}

// Connect starts the MongoDB connection
func (mc *MongoConnection) Connect(onConnectedCallback OnConnectedCallback) {
	mc.onConnectedCallback = onConnectedCallback
	mc.startConnection()
}

// Close closes the MongoDB connection
func (mc *MongoConnection) Close() {
	if mc.client != nil {
		if err := mc.client.Disconnect(context.Background()); err != nil {
			log.Printf("Error closing MongoDB connection: %v\n", err)
		} else {
			log.Println("MongoDB connection closed")
		}
	}
	os.Exit(0)
}

func (mc *MongoConnection) startConnection() {
	log.Printf("Connecting to MongoDB at %s\n", mc.mongoURL)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mc.mongoURL))
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v\n", err)
		return
	}

	mc.client = client
	mc.isConnectedBefore = true
	mc.onConnectedCallback()

	// Start monitoring changes in the server's status
	go mc.monitorServerStatus()
}

func (mc *MongoConnection) monitorServerStatus() {
	// Create a change stream to monitor server status changes
	changeStream, err := mc.client.Watch(context.Background(), mongo.Pipeline{}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		log.Printf("Error creating change stream: %v\n", err)
		return
	}
	defer changeStream.Close(context.Background())

	// Monitor for changes in server status
	for changeStream.Next(context.Background()) {
		var event bson.M
		if err := changeStream.Decode(&event); err != nil {
			log.Printf("Error decoding change stream event: %v\n", err)
			continue
		}

		// Check the operation type to determine the type of event
		operationType, ok := event["operationType"].(string)
		if !ok {
			log.Println("Invalid operation type in change stream event")
			continue
		}

		if operationType == "invalidate" {
			mc.handleDisconnected()
		} else if operationType == "reconnect" {
			mc.handleReconnected()
		}
	}
}

func (mc *MongoConnection) handleReconnected() {
	log.Println("Reconnected to MongoDB")
	mc.onConnectedCallback()
}

func (mc *MongoConnection) handleError(err error) {
	log.Printf("MongoDB connection error: %v\n", err)
}

func (mc *MongoConnection) handleDisconnected() {
	if !mc.isConnectedBefore {
		time.AfterFunc(2*time.Second, func() {
			mc.startConnection()
		})
		log.Println("Retrying MongoDB connection")
	}
}
