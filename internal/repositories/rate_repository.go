package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rate-reader/internal/logger"
	"rate-reader/internal/models"
)

func (rep *MongoRepository) PutRates(ctx context.Context, rates *models.Rates) (*models.Rates, error) {
	log := logger.FromContext(ctx)
	log.Debugf("PutRates %+v", rates)
	if rates.Id.IsZero() {
		ir, err := rep.ratesC.InsertOne(ctx, rates)
		if err != nil {
			log.Errorf("Inserting rates to DB fails: %s", err)
			return nil, err
		}
		id := ir.InsertedID.(primitive.ObjectID)
		rates.Id = id
		return rates, nil
	}
	return nil, fmt.Errorf("rates id is not empty: %s", rates.Id)
}

func (rep *MongoRepository) GetRates(ctx context.Context) (rates []*models.Rates, err error) {
	log := logger.FromContext(ctx)
	log.Debugf("GetRates")

	rates = make([]*models.Rates, 0)
	_, err = rep.ratesC.Find(ctx, &rates)

	if err != nil {
		log.Errorf("GetRates from DB fails: %s", err)
		return nil, err
	}

	return rates, nil
}
