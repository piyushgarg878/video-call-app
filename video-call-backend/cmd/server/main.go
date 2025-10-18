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

	// MongoDB connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://piyushgarg878_db_user:pg878@cluster0.corkolo.mongodb.net/"))
	if err != nil {
		log.Fatal("mongo connect error:", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("mongo disconnect error:", err)
		}
	}()

	db := client.Database("videoapp")

	// Layer setup
	repo := repositories.NewMeetingRepo(db)
	service := services.NewMeetingService(repo)
	meetingHandler := handlers.NewMeetingHandler(service)

	hub := ws.NewHub()
	wsHandler := handlers.NewWSHandler(hub)

	// Routes
	api := app.Group("/api")
	api.Post("/meetings", meetingHandler.Create)
	// Future: api.Get("/meetings/:id", meetingHandler.Get)

	// WebSocket endpoint for signaling
	app.Get("/ws", websocket.New(wsHandler.Handle))

	// Simple health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Video call backend is running 🚀")
	})

	log.Println("Server running on :5000")
	log.Fatal(app.Listen(":5000"))
}