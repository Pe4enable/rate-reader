package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rate-reader/config"
	"rate-reader/logger"
	"rate-reader/models"
	"sync"
	"time"
)

const (
	CRates = "rates"
)

type IRepository interface {
	PutRates(ctx context.Context, rates *models.Rates) (*models.Rates, error)
	GetRates(ctx context.Context) (rates []*models.Rates, err error)
}

type Repository struct {
	Client *mongo.Client
	ratesC *mongo.Collection
	DbName string
}

var (
	rep            IRepository
	onceRepository sync.Once
)

func GetRepository() IRepository {
	onceRepository.Do(func() {
		panic("try to get Repository before it's creation!")
	})
	return rep
}

func New(ctx context.Context, dbConfig config.DB) error {
	var initErr error
	onceRepository.Do(func() {
		rep, initErr = connect(ctx, dbConfig.Host, dbConfig.Name)
	})
	return initErr
}

func connect(ctx context.Context, url, dbName string) (IRepository, error) {
	log := logger.FromContext(ctx)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+url))
	if err != nil {
		log.Errorf("Failed to create connect configuration to MongoDB %s: %v", url, err)
		return nil, err
	}
	for i := 1; i < 6; i++ {
		log.Infof("Attempt %d to connect to MongoDB at %s", i, url)
		err := mongoClient.Ping(ctx, nil)
		if err != nil {
			log.Errorf("Failed to connect to MongoDB at %s: %v", url, err)
			time.Sleep(1 * time.Second)
		} else {
			var db = mongoClient.Database(dbName)
			log.Infof("Connected successfully to MongoDB at %s", url)
			rep = &Repository{
				Client: mongoClient,
				ratesC: db.Collection(CRates),
				DbName: dbName}
			return rep, nil
		}
	}
	return nil, err
}
