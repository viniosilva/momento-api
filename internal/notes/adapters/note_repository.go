package adapters

import (
	"context"
	"errors"
	"time"

	"momento/internal/notes/domain"
	"momento/pkg/listopts"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type noteRepository struct {
	collection *mongo.Collection
}

func NewNoteRepository(collection *mongo.Collection) *noteRepository {
	return &noteRepository{
		collection: collection,
	}
}

func (r *noteRepository) Create(ctx context.Context, note domain.Note) error {
	_, err := r.collection.InsertOne(ctx, note)
	return err
}

func (r *noteRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID, params listopts.ListParams) (listopts.Paginated[domain.Note], error) {
	filter := bson.M{"user_id": userID}

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return listopts.Paginated[domain.Note]{}, err
	}

	findOptions := params.ToFindOptions()
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return listopts.Paginated[domain.Note]{}, err
	}
	defer cursor.Close(ctx)

	var notes []domain.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return listopts.Paginated[domain.Note]{}, err
	}

	return listopts.NewPaginated(notes, totalCount, params.Pagination), nil
}

func (r *noteRepository) GetByIDAndUserID(ctx context.Context, id, userID primitive.ObjectID) (domain.Note, error) {
	filter := bson.M{
		"_id":     id,
		"user_id": userID,
	}

	res := r.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Note{}, domain.ErrNoteNotFound
		}
		return domain.Note{}, err
	}

	var note domain.Note
	if err := res.Decode(&note); err != nil {
		return domain.Note{}, err
	}

	return note, nil
}

func (r *noteRepository) Update(ctx context.Context, note domain.Note) error {
	filter := bson.M{
		"_id":     note.ID,
		"user_id": note.UserID,
	}

	update := bson.M{
		"$set": bson.M{
			"title":      note.Title,
			"content":    note.Content,
			"updated_at": note.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNoteNotFound
	}

	return nil
}

func (r *noteRepository) DeleteByIDAndUserID(ctx context.Context, id, userID primitive.ObjectID) error {
	filter := bson.M{
		"_id":     id,
		"user_id": userID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNoteNotFound
	}

	return nil
}

func (r *noteRepository) ArchiveByIDAndUserID(ctx context.Context, id, userID primitive.ObjectID) error {
	filter := bson.M{
		"_id":         id,
		"user_id":     userID,
		"archived_at": bson.M{"$exists": false},
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
		return domain.ErrNoteNotFound
	}

	return nil
}

func (r *noteRepository) RestoreByIDAndUserID(ctx context.Context, id, userID primitive.ObjectID) error {
	filter := bson.M{
		"_id":         id,
		"user_id":     userID,
		"archived_at": bson.M{"$exists": true},
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
		return domain.ErrNoteNotFound
	}

	return nil
}
