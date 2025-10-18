package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/piyushgarg878/video-call-backend/internal/models"
    "github.com/piyushgarg878/video-call-backend/internal/repositories"
)

type MeetingService struct {
	repo *repository.MeetingRepo
}

func NewMeetingService(repo *repository.MeetingRepo) *MeetingService {
	return &MeetingService{repo: repo}
}

func (s *MeetingService) CreateMeeting(ctx context.Context, title, creatorID string) (*models.Meeting, error) {
	meeting := models.Meeting{
		ID:        uuid.New().String(),
		Title:     title,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateMeeting(ctx, meeting); err != nil {
		return nil, err
	}
	return &meeting, nil
}

func (s *MeetingService) GetMeeting(ctx context.Context, id string) (*models.Meeting, error) {
	return s.repo.GetMeeting(ctx, id)
}