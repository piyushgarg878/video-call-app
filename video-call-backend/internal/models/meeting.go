package models

import "time"


type Meeting struct{
	ID        string    `bson:"_id,omitempty" json:"id"`
	Title     string    `bson:"title" json:"title"`
	CreatorID string    `bson:"creatorId" json:"creatorId"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}