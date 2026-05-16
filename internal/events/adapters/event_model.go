package adapters

import (
	"fmt"
	"log/slog"
	"time"

	"momento/internal/events/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const eventsCollectionName = "events"

type eventMetadataDocument struct {
	ImagePaths []string `bson:"image_paths"`
}

type eventDocument struct {
	ID          primitive.ObjectID     `bson:"_id"`
	OwnerUserID primitive.ObjectID     `bson:"owner_user_id"`
	Title       string                 `bson:"title"`
	Content     string                 `bson:"content"`
	Metadata    *eventMetadataDocument `bson:"metadata,omitempty"`
	CreatedAt   time.Time              `bson:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at"`
	ArchivedAt  *time.Time             `bson:"archived_at"`
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

	doc := eventDocument{
		ID:          id,
		OwnerUserID: ownerUserID,
		Title:       string(e.Title),
		Content:     string(e.Content),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		ArchivedAt:  e.ArchivedAt,
	}

	if e.Metadata != nil {
		imagePaths := make([]string, len(e.Metadata.ImagePaths))
		for i, p := range e.Metadata.ImagePaths {
			imagePaths[i] = string(p)
		}
		doc.Metadata = &eventMetadataDocument{ImagePaths: imagePaths}
	}

	return doc, nil
}

func toEventDomain(d eventDocument) domain.Event {
	evt := domain.Event{
		ID:          d.ID.Hex(),
		OwnerUserID: d.OwnerUserID.Hex(),
		Title:       domain.EventTitle(d.Title),
		Content:     domain.EventContent(d.Content),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		ArchivedAt:  d.ArchivedAt,
	}

	if d.Metadata != nil {
		metadata := domain.NewEventMetadata()
		for _, p := range d.Metadata.ImagePaths {
			path, err := domain.NewImagePath(p)
			if err == nil {
				metadata.ImagePaths = append(metadata.ImagePaths, path)
			} else {
				slog.Warn("skipping invalid image path during deserialization", "path", p, "error", err)
			}
		}
		evt.Metadata = &metadata
	}

	return evt
}

func parseObjectID(hex string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("invalid ID %q: %w", hex, err)
	}
	return id, nil
}
