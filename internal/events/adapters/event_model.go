package adapters

import (
	"fmt"
	"time"

	"momento/internal/events/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const eventsCollectionName = "events"

type eventDocument struct {
	ID          primitive.ObjectID `bson:"_id"`
	OwnerUserID primitive.ObjectID `bson:"owner_user_id"`
	Title       string             `bson:"title"`
	Content     string             `bson:"content"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	ArchivedAt  *time.Time         `bson:"archived_at"`
}

func toEventDocument(e domain.Event) (eventDocument, error) {
	id, err := primitive.ObjectIDFromHex(e.ID)
	if err != nil {
		return eventDocument{}, fmt.Errorf("invalid event ID: %w", err)
	}

	ownerUserID, err := primitive.ObjectIDFromHex(e.OwnerUserID)
	if err != nil {
		return eventDocument{}, fmt.Errorf("invalid user ID: %w", err)
	}

	return eventDocument{
		ID:          id,
		OwnerUserID: ownerUserID,
		Title:       string(e.Title),
		Content:     string(e.Content),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		ArchivedAt:  e.ArchivedAt,
	}, nil
}

func toEventDomain(d eventDocument) domain.Event {
	return domain.Event{
		ID:          d.ID.Hex(),
		OwnerUserID: d.OwnerUserID.Hex(),
		Title:       domain.EventTitle(d.Title),
		Content:     domain.EventContent(d.Content),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		ArchivedAt:  d.ArchivedAt,
	}
}

func parseObjectID(hex string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("invalid ID %q: %w", hex, err)
	}
	return id, nil
}
