package main

import (
	"context"
	"log"
	"time"
	"url-shortener/pkg/shortener"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongo() (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection := client.Database("urlshortener").Collection("url")

	return collection, err
}

func main() {

	collection, err := connectMongo()
	if err != nil {
		log.Fatal(err)
	}

	svc := shortener.NewService(collection)

	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.GET("/:src", svc.Shorten)
	r.GET("/u/:dest", svc.Get)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
