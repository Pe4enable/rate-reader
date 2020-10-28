package rest

import (
	"html/template"
	"net/http"
	"os"

	globalHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func New(handlers *HandlersService) http.Handler {
	r := mux.NewRouter()
	defineRoutes(r, handlers)

	return defineMiddlewares(r)
}

func defineMiddlewares(r *mux.Router) http.Handler {
	router := globalHandlers.LoggingHandler(os.Stdout, r)
	router = globalHandlers.RecoveryHandler()(router)

	return router
}

func defineRoutes(r *mux.Router, handlers *HandlersService) {
	apiSubRouter := r.PathPrefix("/api/v1").Subrouter()
	apiSubRouter.HandleFunc("/help", getHelp).Methods("GET")
	apiSubRouter.HandleFunc("/swagger.yml", swagger).Methods("GET")
	apiSubRouter.HandleFunc("/rates", handlers.GetRates).Methods("GET")

}

func getHelp(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./rest/doc/index.html")
	t.Execute(w, nil)
}

func swagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "./rest/doc/swagger.yml")
}
