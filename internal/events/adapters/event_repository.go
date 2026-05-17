package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"momento/internal/events/domain"
	"momento/pkg/listopts"
)

type eventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *eventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) Create(ctx context.Context, event domain.Event) error {
	row := toEventRow(event)

	query := `INSERT INTO events (id, owner_user_id, title, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query, row.ID, row.OwnerUserID, row.Title, row.Content, row.CreatedAt, row.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert event: %w", err)
	}

	return nil
}

func (r *eventRepository) ListByUserID(ctx context.Context, userID string, params listopts.ListParams) (listopts.Paginated[domain.Event], error) {
	var totalCount int64
	countQuery := `SELECT COUNT(*) FROM events WHERE owner_user_id = $1`
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		return listopts.Paginated[domain.Event]{}, fmt.Errorf("count events: %w", err)
	}

	order := params.ToSQLOrder()
	dataQuery := fmt.Sprintf(`SELECT id, owner_user_id, title, content, created_at, updated_at, archived_at FROM events WHERE owner_user_id = $1 ORDER BY %s LIMIT $2 OFFSET $3`, order)

	limit := int64(params.Pagination.PageSize)
	offset := int64((params.Pagination.Page - 1) * params.Pagination.PageSize)

	var rows []eventRow
	err = r.db.SelectContext(ctx, &rows, dataQuery, userID, limit, offset)
	if err != nil {
		return listopts.Paginated[domain.Event]{}, fmt.Errorf("list events: %w", err)
	}

	if len(rows) == 0 {
		return listopts.NewPaginated([]domain.Event{}, totalCount, params.Pagination), nil
	}

	eventIDs := make([]string, len(rows))
	for i, row := range rows {
		eventIDs[i] = row.ID
	}

	images, err := r.getImagesByEventIDs(ctx, eventIDs)
	if err != nil {
		return listopts.Paginated[domain.Event]{}, fmt.Errorf("get images: %w", err)
	}

	events := make([]domain.Event, len(rows))
	for i, row := range rows {
		events[i] = toEventDomain(row, images[row.ID])
	}

	return listopts.NewPaginated(events, totalCount, params.Pagination), nil
}

func (r *eventRepository) getImagesByEventIDs(ctx context.Context, eventIDs []string) (map[string][]eventImageRow, error) {
	if len(eventIDs) == 0 {
		return nil, nil
	}

	query := `SELECT event_id, path FROM event_images WHERE event_id = ANY($1)`

	var imageRows []eventImageRow
	err := r.db.SelectContext(ctx, &imageRows, query, eventIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]eventImageRow, len(eventIDs))
	for _, img := range imageRows {
		result[img.EventID] = append(result[img.EventID], img)
	}

	return result, nil
}

func (r *eventRepository) getImagesByEventID(ctx context.Context, eventID string) ([]eventImageRow, error) {
	query := `SELECT event_id, path FROM event_images WHERE event_id = $1`

	var rows []eventImageRow
	err := r.db.SelectContext(ctx, &rows, query, eventID)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *eventRepository) GetByIDAndUserID(ctx context.Context, id, userID string) (domain.Event, error) {
	query := `SELECT id, owner_user_id, title, content, created_at, updated_at, archived_at FROM events WHERE id = $1 AND owner_user_id = $2`

	var row eventRow
	err := r.db.GetContext(ctx, &row, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Event{}, domain.ErrEventNotFound
		}

		return domain.Event{}, fmt.Errorf("get event: %w", err)
	}

	images, err := r.getImagesByEventID(ctx, id)
	if err != nil {
		return domain.Event{}, fmt.Errorf("get event images: %w", err)
	}

	return toEventDomain(row, images), nil
}

func (r *eventRepository) Update(ctx context.Context, event domain.Event) error {
	query := `UPDATE events SET title = $1, content = $2, updated_at = $3 WHERE id = $4 AND owner_user_id = $5`

	result, err := r.db.ExecContext(ctx, query, string(event.Title), string(event.Content), event.UpdatedAt, event.ID, event.OwnerUserID)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) DeleteByIDAndUserID(ctx context.Context, id, userID string) error {
	query := `DELETE FROM events WHERE id = $1 AND owner_user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) ArchiveByIDAndUserID(ctx context.Context, id, userID string) error {
	now := time.Now().UTC()
	query := `UPDATE events SET archived_at = $1, updated_at = $2 WHERE id = $3 AND owner_user_id = $4 AND archived_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, now, now, id, userID)
	if err != nil {
		return fmt.Errorf("archive event: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) RestoreByIDAndUserID(ctx context.Context, id, userID string) error {
	now := time.Now().UTC()
	query := `UPDATE events SET archived_at = NULL, updated_at = $1 WHERE id = $2 AND owner_user_id = $3 AND archived_at IS NOT NULL`

	result, err := r.db.ExecContext(ctx, query, now, id, userID)
	if err != nil {
		return fmt.Errorf("restore event: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) AddImage(ctx context.Context, eventID, userID string, path domain.ImagePath) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var count int
	countQuery := `SELECT COUNT(*) FROM event_images WHERE event_id = $1`
	if err := tx.QueryRowContext(ctx, countQuery, eventID).Scan(&count); err != nil {
		return fmt.Errorf("count images: %w", err)
	}

	if count >= domain.MaxImages {
		return domain.ErrMaxImagesReached
	}

	insertQuery := `INSERT INTO event_images (event_id, path) VALUES ($1, $2)`
	if _, err := tx.ExecContext(ctx, insertQuery, eventID, string(path)); err != nil {
		if isUniqueViolation(err) {
			return nil
		}

		return fmt.Errorf("insert image: %w", err)
	}

	updateQuery := `UPDATE events SET updated_at = $1 WHERE id = $2 AND owner_user_id = $3`
	result, err := tx.ExecContext(ctx, updateQuery, time.Now().UTC(), eventID, userID)
	if err != nil {
		return fmt.Errorf("update event timestamp: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrEventNotFound
	}

	return tx.Commit()
}

func (r *eventRepository) RemoveImage(ctx context.Context, eventID, userID string, path domain.ImagePath) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	deleteQuery := `DELETE FROM event_images WHERE event_id = $1 AND path = $2`
	result, err := tx.ExecContext(ctx, deleteQuery, eventID, string(path))
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return domain.ErrImageNotFound
	}

	updateQuery := `UPDATE events SET updated_at = $1 WHERE id = $2 AND owner_user_id = $3`
	_, err = tx.ExecContext(ctx, updateQuery, time.Now().UTC(), eventID, userID)
	if err != nil {
		return fmt.Errorf("update event timestamp: %w", err)
	}

	return tx.Commit()
}

func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate key") || contains(err.Error(), "unique constraint") || contains(err.Error(), "23505"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
