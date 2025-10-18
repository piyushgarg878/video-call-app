package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/piyushgarg878/video-call-backend/internal/handlers"
	"github.com/piyushgarg878/video-call-backend/internal/repositories"
	"github.com/piyushgarg878/video-call-backend/internal/services"
	"github.com/piyushgarg878/video-call-backend/internal/ws"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func main() {
	app := fiber.New()

	// MongoDB
	ctx:=context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db := client.Database("videoapp")
	if err!=nil{
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	// Layers
	repo := repository.NewMeetingRepo(db)
	service := services.NewMeetingService(repo)
	meetingHandler := handlers.NewMeetingHandler(service)

	hub := ws.NewHub()
	wsHandler := handlers.NewWSHandler(hub)

	// Routes
	api := app.Group("/api")
	api.Post("/meetings", meetingHandler.Create)

	app.Get("/ws", websocket.New(wsHandler.Handle))

	log.Println("Server running on :5000")
	log.Fatal(app.Listen(":5000"))
}