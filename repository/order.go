package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/srrmendez/private-api-order/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRepository(mongoClient *mongo.Client, database string, collection string) *repository {
	dc := mongoClient.Database(database).Collection(collection)

	return &repository{
		collection: dc,
	}
}

func (r *repository) All(ctx context.Context, appID string) ([]model.Order, error) {
	findOptions := options.FindOptions{}
	findOptions.SetSort(bson.D{{"created_at", -1}})

	cursor, err := r.collection.Find(ctx, bson.D{{"app_id", appID}}, &findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	orders := make([]model.Order, 0)

	for cursor.Next(ctx) {
		var order model.Order

		err = cursor.Decode(&order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *repository) Upsert(ctx context.Context, order model.Order) (*model.Order, error) {
	now := time.Now().Unix()

	if order.ID == "" {
		order.ID = uuid.NewString()
		order.CreatedAt = model.CustomTimeStamp(now)
	}

	order.UpdatedAt = model.CustomTimeStamp(now)

	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *repository) Search(ctx context.Context, appID string, category *model.CategoryType, status *model.OrderStatusType, orderType *model.OrderType) ([]model.Order, error) {
	if status == nil && orderType == nil && category != nil {
		return nil, errors.New("must provide a filter value")
	}

	findOptions := options.FindOptions{}
	findOptions.SetSort(bson.D{{"created_at", -1}})

	filter := bson.D{{"app_id", appID}}

	if category != nil {
		filter = append(filter, bson.E{"category", *category})
	}

	if status != nil {
		filter = append(filter, bson.E{"status", *status})
	}

	if orderType != nil {
		filter = append(filter, bson.E{"type", *orderType})
	}

	cursor, err := r.collection.Find(ctx, filter, &findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	orders := make([]model.Order, 0)

	for cursor.Next(ctx) {
		var order model.Order

		err = cursor.Decode(&order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *repository) Get(ctx context.Context, id string) (*model.Order, error) {
	filter := bson.D{{"_id", id}}

	var order model.Order

	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &order, nil
}

func (r *repository) GetByExternalID(ctx context.Context, externalID string) (*model.Order, error) {
	filter := bson.D{{"external_id", externalID}}

	var order model.Order

	err := r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &order, nil
}
