package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/srrmendez/private-api-offers/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRepository(client *mongo.Client, database string, table string) *repository {
	return &repository{
		collection: client.Database(database).Collection(table),
	}
}

func (r *repository) All(ctx context.Context) ([]model.Offer, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	offers := make([]model.Offer, 0)

	for cursor.Next(ctx) {
		var offer model.Offer

		err = cursor.Decode(&offer)
		if err != nil {
			return nil, err
		}

		offers = append(offers, offer)
	}

	return offers, nil
}

func (r *repository) Upsert(ctx context.Context, offer model.Offer) (*model.Offer, error) {
	r.initializeOffer(&offer)

	upsert := true

	_, err := r.collection.UpdateOne(ctx, bson.D{{"_id", offer.ID}},
		bson.D{{"$set", offer}}, &options.UpdateOptions{
			Upsert: &upsert,
		})
	if err != nil {
		return nil, err
	}

	return &offer, nil
}

func (r *repository) initializeOffer(offer *model.Offer) {
	now := time.Now().Unix()

	if offer.ID == "" {
		offer.ID = uuid.NewString()
		offer.CreatedAt = model.CustomTimeStamp(now)
	}

	offer.UpdatedAt = model.CustomTimeStamp(now)
}

func (r *repository) Get(ctx context.Context, id string) (*model.Offer, error) {
	filter := bson.D{{"_id", id}}

	var offer model.Offer

	err := r.collection.FindOne(ctx, filter).Decode(&offer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &offer, nil
}

func (r *repository) GetByExternalID(ctx context.Context, id string) (*model.Offer, error) {
	filter := bson.D{{"external_id", id}}

	var offer model.Offer

	err := r.collection.FindOne(ctx, filter).Decode(&offer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &offer, nil
}

func (r *repository) Search(ctx context.Context, active *bool, category *model.CategoryType) ([]model.Offer, error) {
	now := time.Now().Unix()

	query := bson.D{}

	if active != nil {
		query = bson.D{
			{"effective_date", bson.D{{"$lte", now}}},
			{"expiration_date", bson.D{{"$gte", now}}},
		}

		if !*active {
			query = bson.D{
				{"$or", []bson.D{
					{{"effective_date", bson.D{{"$gt", now}}}},
					{{"expiration_date", bson.D{{"$lt", now}}}},
				}},
			}
		}
	}

	if category != nil {
		query = append(query, bson.D{{"category", *category}}...)
	}

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	offers := make([]model.Offer, 0)

	for cursor.Next(ctx) {
		var offer model.Offer

		err = cursor.Decode(&offer)
		if err != nil {
			return nil, err
		}

		offers = append(offers, offer)
	}

	return offers, nil
}

func (r *repository) RemoveByExternalID(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.D{{"external_id", id}})
	if err != nil {
		return err
	}

	return nil
}
