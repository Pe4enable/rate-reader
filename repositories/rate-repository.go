package repositories

import (
	"context"
	"fmt"
	"github.com/rate-reader/logger"
	"github.com/rate-reader/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (rep *Repository) PutRates(ctx context.Context, rates *models.Rates) (*models.Rates, error) {
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
