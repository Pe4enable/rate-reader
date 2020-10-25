package handlers

import (
	"net/http"
	"rate-reader/logger"
	"rate-reader/repositories"
)

type HandlersService struct {
	rep repositories.IRepository
}

func NewHandlerService(
	rep repositories.IRepository) *HandlersService {
	return &HandlersService{
		rep: rep,
	}
}

func (s *HandlersService) GetRates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	log.Infof("Rest GetRates")

	rates, err := s.rep.GetRates(ctx)
	if err != nil {
		log.Errorf("Error during get rates from DB: %s.", err)
		jsonErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, rates, http.StatusOK)
}
