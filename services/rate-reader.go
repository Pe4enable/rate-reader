package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rate-reader/config"
	"github.com/rate-reader/logger"
	"github.com/rate-reader/models"
	"github.com/rate-reader/repositories"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	requestDelay = 10000
	getMarketSummaries = "getmarketsummaries"
	delayAfterError = time.Second * 30
)

type IRateReader interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context)
}

type rateReader struct {
	rp               repositories.IRepository
	conf             *config.Config
	delay 			 time.Duration
	stopListen       chan struct{}
	httpClient http.Client
}

var (
	reader         IRateReader
	onceRateReader sync.Once
)

func NewReader(ctx context.Context, config *config.Config, rp repositories.IRepository) error {
	onceRateReader.Do(func() {

		reader = &rateReader{
			rp:               rp,
			conf:             config,
			delay: time.Duration(config.Delay) * time.Millisecond,
			httpClient: http.Client{
				Timeout:   time.Duration(requestDelay) * time.Millisecond,
				Transport: &http.Transport{},
			},
		}
	})
	return nil
}

// GetRateReader returns rates reader instance.
// Client must be previously created with New(), in another case function throws panic
func GetRateReader() IRateReader {
	onceRateReader.Do(func() {
		panic("try to get rate reader before it's creation!")
	})
	return reader
}

func (rr *rateReader) Start(ctx context.Context) (err error) {
	log := logger.FromContext(ctx)

	log.Infof("Start reading rates from %s", rr.conf.Source)
	go func() {
		for {
			select {
			case <-rr.stopListen:
				log.Infof("Stop reading rates.")
				return
			default:
			}

			raw, err := rr.restRequest(ctx, http.MethodGet, rr.conf.Source + getMarketSummaries, nil)
			if err != nil {
				log.Errorf("failed to get rates from source: %s", err)
				time.Sleep(delayAfterError)
				continue
			}

			currentRates := &models.Rates{}
			if err := json.Unmarshal(raw, currentRates); err != nil {
				log.Errorf("failed to parse current rates: %s", err)
				time.Sleep(delayAfterError)
				continue
			}

			currentRates.TimeStamp =time.Now()
			currentRates, err = rr.rp.PutRates(ctx, currentRates)
			if err != nil {
				log.Errorf("Put rates to db error: %s", err)
			}
			time.Sleep(rr.delay)
		}
	}()

	return
}

func (rr *rateReader) Stop(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Info("Stop reading rates.")
	rr.stopListen <- struct{}{}
	return
}

func (rr *rateReader) restRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := rr.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("wrong response status %d, body %s", resp.StatusCode, string(body))

	}
	respBody, err := ioutil.ReadAll(resp.Body)
	log := logger.FromContext(ctx)
	log.Debugf("response: %s", string(respBody))
	return respBody, err
}