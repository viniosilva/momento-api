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

func NewNoteRepository(db *mongo.Database) *noteRepository {
	return &noteRepository{
		collection: db.Collection(notesCollectionName),
	}
}

func (r *noteRepository) Create(ctx context.Context, note domain.Note) error {
	doc, err := toNoteDocument(note)
	if err != nil {
		return err
	}

	_, err = r.collection.InsertOne(ctx, doc)
	return err
}

func (r *noteRepository) ListByUserID(ctx context.Context, userID string, params listopts.ListParams) (listopts.Paginated[domain.Note], error) {
	uid, err := parseObjectID(userID)
	if err != nil {
		return listopts.Paginated[domain.Note]{}, err
	}

	filter := bson.M{"user_id": uid}

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

	var docs []noteDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return listopts.Paginated[domain.Note]{}, err
	}

	notes := make([]domain.Note, len(docs))
	for i, d := range docs {
		notes[i] = toNoteDomain(d)
	}

	return listopts.NewPaginated(notes, totalCount, params.Pagination), nil
}

func (r *noteRepository) GetByIDAndUserID(ctx context.Context, id, userID string) (domain.Note, error) {
	oid, err := parseObjectID(id)
	if err != nil {
		return domain.Note{}, err
	}

	uid, err := parseObjectID(userID)
	if err != nil {
		return domain.Note{}, err
	}

	filter := bson.M{
		"_id":     oid,
		"user_id": uid,
	}

	res := r.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Note{}, domain.ErrNoteNotFound
		}
		return domain.Note{}, err
	}

	var doc noteDocument
	if err := res.Decode(&doc); err != nil {
		return domain.Note{}, err
	}

	return toNoteDomain(doc), nil
}

func (r *noteRepository) Update(ctx context.Context, note domain.Note) error {
	oid, err := parseObjectID(note.ID)
	if err != nil {
		return err
	}

	uid, err := parseObjectID(note.UserID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":     oid,
		"user_id": uid,
	}

	update := bson.M{
		"$set": bson.M{
			"title":      string(note.Title),
			"content":    string(note.Content),
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

func (r *noteRepository) DeleteByIDAndUserID(ctx context.Context, id, userID string) error {
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
		return domain.ErrNoteNotFound
	}

	return nil
}

func (r *noteRepository) ArchiveByIDAndUserID(ctx context.Context, id, userID string) error {
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
		return domain.ErrNoteNotFound
	}

	return nil
}

func (r *noteRepository) RestoreByIDAndUserID(ctx context.Context, id, userID string) error {
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
		return domain.ErrNoteNotFound
	}

	return nil
}
