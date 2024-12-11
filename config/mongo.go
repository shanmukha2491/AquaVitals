package config

import (
	"context"
	"errors"
	"fmt"
	"log"

	"time"


	"github.com/shanmukha2491/AquaVitals/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

var _UserCollection *mongo.Collection

var _AdminCollection_ *mongo.Collection

func ConnectDB() *mongo.Client {
	// Set MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb+srv://shanmukha2491:AquaVitals@cluster0.4v93m.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	// Create a context with a timeout to avoid indefinite hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error occurred while connecting to database: %v", err)
		return nil
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil
	}

	log.Println("Database connection was successful!")
	Client = client
	return client
}

func UserCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database("AquaVitalsStorage").Collection("User")
	_UserCollection = collection
	return collection
}

func AdminCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database("AquaVitalsStorage").Collection("admin")
	_AdminCollection_ = collection
	return collection
}

func CreateUser(user *model.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := _UserCollection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error:", err.Error())
		return errors.New("failed to create user: " + err.Error())
	}
	return nil

}

func FindOne(email string, userName string) (model.User, error) {
	filter := bson.M{"email": email, "username": userName}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user model.User
	err := _UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func RegisterSensor(newSensor model.Sensor, userName string) error {
	// MongoDB filter to find the user document by username
	filter := bson.M{"username": userName}


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()


	// MongoDB update to push the new sensor to the `sensors` array
	pushUpdate := bson.M{
		"$push": bson.M{"sensors": newSensor},
	}

	// Perform the push update
	result, err := _UserCollection.UpdateOne(ctx, filter, pushUpdate)
	if err != nil {
		return fmt.Errorf("failed to push new sensor: %v", err)
	}

	// Log the result
	fmt.Printf("Matched %d documents and modified %d documents.\n", result.MatchedCount, result.ModifiedCount)

	return nil
}

