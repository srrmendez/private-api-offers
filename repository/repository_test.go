package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/srrmendez/private-api-order/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestRepository(t *testing.T) {
	t.Run("successfully insert order", func(t *testing.T) {
		t.Skip()
		// Given
		externalID := uuid.NewString()

		order := model.Order{
			Status:     model.CompleteOrderStatus,
			Type:       model.VPSOrder,
			ExternalID: &externalID,
		}

		ctx := context.Background()

		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		require.NoError(t, err)

		repo := NewRepository(mongoClient, "operations", "operations")

		// When
		o, err := repo.Upsert(ctx, order)

		// Then
		require.NoError(t, err)
		require.NotNil(t, o)
		assert.NotEqual(t, "", o.ID)
	})

	t.Run("successfully list all orders", func(t *testing.T) {
		t.Skip()
		// Given
		ctx := context.Background()

		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		require.NoError(t, err)

		repo := NewRepository(mongoClient, "operations", "operations")

		// When
		o, err := repo.All(ctx, "")

		// Then
		require.NoError(t, err)
		require.NotNil(t, o)
		assert.Greater(t, 0, len(o))
	})

	t.Run("successfully get order", func(t *testing.T) {
		t.Skip()
		// Given
		ctx := context.Background()

		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		require.NoError(t, err)

		repo := NewRepository(mongoClient, "operations", "operations")

		// When
		o, err := repo.Get(ctx, "a6750743-be51-4742-a278-9d74bfa59a9f")

		// Then
		require.NoError(t, err)
		require.NotNil(t, o)
	})

	t.Run("failed get order, not found", func(t *testing.T) {
		t.Skip()
		// Given
		ctx := context.Background()

		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		require.NoError(t, err)

		repo := NewRepository(mongoClient, "operations", "operations")

		// When
		o, err := repo.Get(ctx, "fake")

		// Then
		require.NoError(t, err)
		require.Nil(t, o)
	})

	t.Run("successfully list all completed orders", func(t *testing.T) {
		t.Skip()
		// Given
		ctx := context.Background()

		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		require.NoError(t, err)

		repo := NewRepository(mongoClient, "operations", "operations")

		status := model.CompleteOrderStatus

		// When
		o, err := repo.Search(ctx, "", nil, &status, nil)

		// Then
		require.NoError(t, err)
		require.Nil(t, o)
	})

	t.Run("successfully list all completed orders with VPS type", func(t *testing.T) {
		// Given
		ctx := context.Background()

		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		require.NoError(t, err)

		repo := NewRepository(mongoClient, "operations", "operations")

		status := model.CompleteOrderStatus
		tp := model.VPSOrder

		// When
		o, err := repo.Search(ctx, "", nil, &status, &tp)

		// Then
		require.NoError(t, err)
		require.Nil(t, o)
	})

}
