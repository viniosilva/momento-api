package adapters

import (
	"time"

	"momento/internal/events/domain"
)

type eventRow struct {
	ID          string     `db:"id"`
	OwnerUserID string     `db:"owner_user_id"`
	Title       string     `db:"title"`
	Content     string     `db:"content"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	ArchivedAt  *time.Time `db:"archived_at"`
}

type eventImageRow struct {
	EventID string `db:"event_id"`
	Path    string `db:"path"`
}

func toEventRow(e domain.Event) eventRow {
	return eventRow{
		ID:          e.ID,
		OwnerUserID: e.OwnerUserID,
		Title:       string(e.Title),
		Content:     string(e.Content),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		ArchivedAt:  e.ArchivedAt,
	}
}

func toEventDomain(row eventRow, imageRows []eventImageRow) domain.Event {
	evt := domain.Event{
		ID:          row.ID,
		OwnerUserID: row.OwnerUserID,
		Title:       domain.EventTitle(row.Title),
		Content:     domain.EventContent(row.Content),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		ArchivedAt:  row.ArchivedAt,
	}

	if len(imageRows) > 0 {
		metadata := domain.NewEventMetadata()
		for _, img := range imageRows {
			path, err := domain.NewImagePath(img.Path)
			if err == nil {
				metadata.AddImage(path)
			}
		}
		evt.Metadata = &metadata
	}

	return evt
}
