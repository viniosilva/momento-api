package adapters

import (
	"context"
	"errors"
	"time"

	"momento/internal/events/domain"
	"momento/pkg/listopts"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type eventRepository struct {
	collection *mongo.Collection
}

func NewEventRepository(db *mongo.Database) *eventRepository {
	return &eventRepository{
		collection: db.Collection(eventsCollectionName),
	}
}

func (r *eventRepository) Create(ctx context.Context, event domain.Event) error {
	doc, err := toEventDocument(event)
	if err != nil {
		return err
	}

	_, err = r.collection.InsertOne(ctx, doc)
	return err
}

func (r *eventRepository) ListByUserID(ctx context.Context, userID string, params listopts.ListParams) (listopts.Paginated[domain.Event], error) {
	uid, err := parseObjectID(userID)
	if err != nil {
		return listopts.Paginated[domain.Event]{}, err
	}

	filter := bson.M{"user_id": uid}

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return listopts.Paginated[domain.Event]{}, err
	}

	findOptions := params.ToFindOptions()
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return listopts.Paginated[domain.Event]{}, err
	}
	defer cursor.Close(ctx)

	var docs []eventDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return listopts.Paginated[domain.Event]{}, err
	}

	events := make([]domain.Event, len(docs))
	for i, d := range docs {
		events[i] = toEventDomain(d)
	}

	return listopts.NewPaginated(events, totalCount, params.Pagination), nil
}

func (r *eventRepository) GetByIDAndUserID(ctx context.Context, id, userID string) (domain.Event, error) {
	oid, err := parseObjectID(id)
	if err != nil {
		return domain.Event{}, err
	}

	uid, err := parseObjectID(userID)
	if err != nil {
		return domain.Event{}, err
	}

	filter := bson.M{
		"_id":     oid,
		"user_id": uid,
	}

	res := r.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Event{}, domain.ErrEventNotFound
		}
		return domain.Event{}, err
	}

	var doc eventDocument
	if err := res.Decode(&doc); err != nil {
		return domain.Event{}, err
	}

	return toEventDomain(doc), nil
}

func (r *eventRepository) Update(ctx context.Context, event domain.Event) error {
	oid, err := parseObjectID(event.ID)
	if err != nil {
		return err
	}

	uid, err := parseObjectID(event.OwnerUserID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":     oid,
		"user_id": uid,
	}

	update := bson.M{
		"$set": bson.M{
			"title":      string(event.Title),
			"content":    string(event.Content),
			"updated_at": event.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) DeleteByIDAndUserID(ctx context.Context, id, userID string) error {
	oid, err := parseObjectID(id)
	if err != nil {
		return err
	}

	uid, err := parseObjectID(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":     oid,
		"user_id": uid,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) ArchiveByIDAndUserID(ctx context.Context, id, userID string) error {
	oid, err := parseObjectID(id)
	if err != nil {
		return err
	}

	uid, err := parseObjectID(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":         oid,
		"user_id":     uid,
		"archived_at": nil,
	}

	update := bson.M{
		"$set": bson.M{
			"archived_at": primitive.NewDateTimeFromTime(time.Now().UTC()),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventRepository) RestoreByIDAndUserID(ctx context.Context, id, userID string) error {
	oid, err := parseObjectID(id)
	if err != nil {
		return err
	}

	uid, err := parseObjectID(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":         oid,
		"user_id":     uid,
		"archived_at": bson.M{"$ne": nil},
	}

	update := bson.M{
		"$set": bson.M{
			"archived_at": nil,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}
