package config

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/shanmukha2491/AquaVitals/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

var _UserCollection *mongo.Collection

var _AdminCollection_ *mongo.Collection

func ConnectDB() *mongo.Client {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("error loading url", err)
		return nil
	}

	MongoURl := os.Getenv("MONGO_URL")

	clientOptions := options.Client().ApplyURI(MongoURl)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error Occured while connecting to database")
		return nil
	}
	Client = client
	log.Println("DataBase connection is Success")
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
