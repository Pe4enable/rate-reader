package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rate-reader/internal/config"
	"rate-reader/internal/logger"
	"rate-reader/internal/models"
	"time"
)

const (
	CRates = "rates"
)

type Repository interface {
	PutRates(ctx context.Context, rates *models.Rates) (*models.Rates, error)
	GetRates(ctx context.Context) (rates []*models.Rates, err error)
}

type MongoRepository struct {
	Client *mongo.Client
	ratesC *mongo.Collection
}

func New(ctx context.Context, dbConfig config.DB) (*MongoRepository, error) {
	log := logger.FromContext(ctx)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+dbConfig.Host))
	if err != nil {
		log.Errorf("Failed to create connect configuration to MongoDB %s: %v", dbConfig.Host, err)
		return nil, err
	}
	for i := 1; i < 6; i++ {
		log.Infof("Attempt %d to connect to MongoDB at %s", i, dbConfig.Host)
		err := mongoClient.Ping(ctx, nil)
		if err != nil {
			log.Errorf("Failed to connect to MongoDB at %s: %v", dbConfig.Host, err)
			time.Sleep(1 * time.Second)
		} else {
			var db = mongoClient.Database(dbConfig.Name)
			log.Infof("Connected successfully to MongoDB at %s", dbConfig.Host)
			rep := &MongoRepository{
				Client: mongoClient,
				ratesC: db.Collection(CRates)}
			return rep, nil
		}
	}
	return nil, err
}
