package repository

import (
	"context"
	"testing"
	"time"

	"github.com/srrmendez/private-api-offers/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestOfferRepository(t *testing.T) {
	t.Run("succssfully create offer", func(t *testing.T) {
		t.Skip()
		// Given
		ctx := context.Background()

		mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		repo := NewRepository(mongoClient, "service_layer", "offers")

		effDate := time.Now().Add(-72 * time.Hour).Unix()
		expirDate := time.Now().Add(72 * time.Hour).Unix()

		offer := model.Offer{
			Name:           "offer 4",
			ClientType:     model.IndividualClienType,
			Paymode:        model.PrepaidPayMode,
			MonthlyFee:     200,
			OneOfFee:       10,
			Metadata:       map[string]string{"t1": "v"},
			EffectiveDate:  model.CustomTimeStamp(effDate),
			ExpirationDate: model.CustomTimeStamp(expirDate),
			Category:       "datacenter",
		}

		// When
		of, err := repo.Upsert(ctx, offer)

		// Then
		require.NoError(t, err)
		require.NotNil(t, of)
	})

	t.Run("successfully create offers batch", func(t *testing.T) {
		t.Skip()
		// Given
		ctx := context.Background()

		mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		repo := NewRepository(mongoClient, "service_layer", "offers")

		effDate := time.Now().Add(1 * time.Second).Unix()
		expirDate := time.Now().Add(72 * time.Hour).Unix()

		description := "description 1"

		offer := []model.Offer{
			{
				Name:           "offer 2",
				ClientType:     model.IndividualClienType,
				Paymode:        model.PrepaidPayMode,
				StandAlone:     true,
				MonthlyFee:     100,
				OneOfFee:       10,
				Metadata:       map[string]string{"t": "v"},
				EffectiveDate:  model.CustomTimeStamp(effDate),
				ExpirationDate: model.CustomTimeStamp(expirDate),
				Category:       "data center",
			},
			{
				Name:        "offer 3",
				ClientType:  model.IndividualClienType,
				Paymode:     model.PrepaidPayMode,
				Description: &description,
				Childrens: []model.Offer{{
					Name:           "offer 4",
					ClientType:     model.IndividualClienType,
					Paymode:        model.PrepaidPayMode,
					StandAlone:     true,
					MonthlyFee:     100,
					OneOfFee:       10,
					Metadata:       map[string]string{"t": "v"},
					EffectiveDate:  model.CustomTimeStamp(effDate),
					ExpirationDate: model.CustomTimeStamp(expirDate),
					Category:       "yellow pages",
				}},
				StandAlone:     true,
				MonthlyFee:     200,
				OneOfFee:       40,
				Metadata:       map[string]string{"t": "v"},
				EffectiveDate:  model.CustomTimeStamp(effDate),
				ExpirationDate: model.CustomTimeStamp(expirDate),
				Category:       "data center",
			},
		}

		// When
		err := repo.BatchUpsert(ctx, offer)

		// Then
		require.NoError(t, err)
	})

	t.Run("succssfully get offer", func(t *testing.T) {
		t.Skip()
		// Given
		ctx := context.Background()

		mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		repo := NewRepository(mongoClient, "service_layer", "offers")

		// When
		of, err := repo.Get(ctx, "a7fca57a-86bc-49f0-b68a-60f52d932d11")

		// Then
		require.NoError(t, err)
		require.NotNil(t, of)
	})

	t.Run("succssfully search active offers", func(t *testing.T) {
		t.Skip()
		// Given
		ctx := context.Background()

		mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		repo := NewRepository(mongoClient, "service_layer", "offers")

		// When
		offers, err := repo.Search(ctx, true)

		// Then
		require.NoError(t, err)
		require.NotNil(t, offers)
		assert.Equal(t, 4, len(offers))
	})

	t.Run("succssfully search inactive offers", func(t *testing.T) {
		t.Skip()
		// Given
		ctx := context.Background()

		mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		repo := NewRepository(mongoClient, "service_layer", "offers")

		// When
		offers, err := repo.Search(ctx, false)

		// Then
		require.NoError(t, err)
		require.NotNil(t, offers)
		assert.Equal(t, 4, len(offers))
	})

}
