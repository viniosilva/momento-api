package infrastructure

import (
	"context"

	"pinnado/internal/notes/domain"
	shareddto "pinnado/internal/shared/application/dto"
	sharedinfra "pinnado/internal/shared/infrastructure"

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

func (r *noteRepository) ListByUserID(ctx context.Context, userID primitive.ObjectID, params sharedinfra.ListParams) (shareddto.Paginated[domain.Note], error) {
	filter := bson.M{"user_id": userID}

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return shareddto.Paginated[domain.Note]{}, err
	}

	findOptions := params.ToFindOptions()
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return shareddto.Paginated[domain.Note]{}, err
	}
	defer cursor.Close(ctx)

	var notes []domain.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return shareddto.Paginated[domain.Note]{}, err
	}

	return shareddto.NewPaginated(notes, totalCount, params.Pagination), nil
}
