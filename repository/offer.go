package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/srrmendez/private-api-offers/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	_, err := r.collection.InsertOne(ctx, offer)
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

func (r *repository) BatchUpsert(ctx context.Context, offers []model.Offer) error {
	ofDocs := make([]interface{}, 0, len(offers))

	for _, offer := range offers {
		r.initializeOffer(&offer)

		offDoc := bson.D{
			{"_id", offer.ID},
			{"created_at", offer.CreatedAt},
			{"updated_at", offer.UpdatedAt},
			{"name", offer.Name},
			{"code", offer.Code},
			{"client_type", offer.ClientType},
			{"pay_mode", offer.Paymode},
			{"standalone", offer.StandAlone},
			{"category", offer.Category},
			{"effective_date", offer.EffectiveDate},
			{"expiration_date", offer.ExpirationDate},
			{"monthly_fee", offer.MonthlyFee},
			{"one_of_fee", offer.OneOfFee},
		}

		if offer.ExternalID != nil {
			offDoc = append(offDoc, bson.E{"external_id", offer.ExternalID})
		}

		if offer.Description != nil {
			offDoc = append(offDoc, bson.E{"description", offer.Description})
		}

		if len(offer.Childrens) > 0 {
			for i := range offer.Childrens {
				r.initializeOffer(&offer.Childrens[i])
			}

			offDoc = append(offDoc, bson.E{"childrens", offer.Childrens})
		}

		if offer.Code != nil {
			offDoc = append(offDoc, bson.E{"code", offer.Code})
		}

		if offer.Metadata != nil {
			offDoc = append(offDoc, bson.E{"metadata", offer.Metadata})
		}

		ofDocs = append(ofDocs, offDoc)
	}

	_, err := r.collection.InsertMany(ctx, ofDocs)
	if err != nil {
		return err
	}
	return nil
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

func (r *repository) Search(ctx context.Context, active bool) ([]model.Offer, error) {
	now := time.Now().Unix()

	query := bson.D{
		{"effective_date", bson.D{{"$lte", now}}},
		{"expiration_date", bson.D{{"$gte", now}}},
	}

	if !active {
		query = bson.D{
			{"$or", []bson.D{
				{{"effective_date", bson.D{{"$gt", now}}}},
				{{"expiration_date", bson.D{{"$lt", now}}}},
			}},
		}
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

func (r *repository) Remove(ctx context.Context, ids []string) error {
	query := bson.D{
		{"_id", bson.D{{"$in", ids}}},
	}

	d, err := r.collection.DeleteMany(ctx, query)
	if err != nil {
		return err
	}

	if d.DeletedCount == 0 {
		return errors.New("ids to remove not found")
	}

	return nil
}
