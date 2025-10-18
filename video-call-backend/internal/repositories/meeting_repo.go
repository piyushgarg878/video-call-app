package repository

import (
	"context"

	"github.com/piyushgarg878/video-call-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MeetingRepo struct {
	col *mongo.Collection
}

func NewMeetingRepo(db *mongo.Database) *MeetingRepo {
	return &MeetingRepo{col: db.Collection("meetings")}
}

func (r *MeetingRepo) CreateMeeting(ctx context.Context, meeting models.Meeting) error {
	_, err := r.col.InsertOne(ctx, meeting)
	return err
}

func (r *MeetingRepo) GetMeeting(ctx context.Context, id string) (*models.Meeting, error) {
	var m models.Meeting
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	return &m, err
}