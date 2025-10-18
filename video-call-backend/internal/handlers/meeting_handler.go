package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/piyushgarg878/video-call-backend/internal/services"
)

type MeetingHandler struct {
	svc *services.MeetingService
}

func NewMeetingHandler(s *services.MeetingService) *MeetingHandler {
	return &MeetingHandler{s}
}

func (h *MeetingHandler) Create(c *fiber.Ctx) error {
	var req struct { Title string `json:"title"` }
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	meeting, err := h.svc.CreateMeeting(c.Context(), req.Title, "demo-user")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(meeting)
}