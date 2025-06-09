package repository

import (
	"context"
	"fmt"

	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/pkg/database"
	"go.elastic.co/apm"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertNewMessage(ctx context.Context, data models.MessagePayload) error {
	// Insert the message into the MongoDB collection
	_, err := database.MongoDB.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func GetMessageHistory(ctx context.Context) ([]models.MessagePayload, error) {
	span, _ := apm.StartSpan(ctx, "GetMessageHistory", "repository")
	defer span.End()
	var (
		err  error
		resp []models.MessagePayload
	)
	cursor, err := database.MongoDB.Find(ctx, bson.D{})
	if err != nil {
		return resp, fmt.Errorf("failed to find messages: %v", err)
	}
	for cursor.Next(ctx) {
		payload := models.MessagePayload{}
		if err := cursor.Decode(&payload); err != nil {
			return resp, fmt.Errorf("failed to decode message: %v", err)
		}
		resp = append(resp, payload)
	}
	return resp, nil
}
