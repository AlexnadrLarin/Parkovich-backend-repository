package internal

import (
    "net/http"
	
    "go-parkovich/microservices/analytics/internal/api"
    "go-parkovich/microservices/analytics/internal/database"

    "github.com/gorilla/mux"
)

func SetupRouter(repo *database.UserEventsRepository) *mux.Router {
    r := mux.NewRouter()

    r.HandleFunc("/user-action", api.HandleUserAction(repo)).Methods(http.MethodPost)

    r.HandleFunc("/user-actions", api.GetAllUserActions(repo)).Methods(http.MethodGet)

    r.HandleFunc("/user-actions/{user_id}", api.GetUserActions(repo)).Methods(http.MethodGet)

    return r
}
