package internal

import (
    "go-parkovich/microservices/analytics/internal/api"
    "go-parkovich/microservices/analytics/internal/database"
    "go-parkovich/microservices/analytics/internal/middleware"

    "github.com/gorilla/mux"

)

func SetupRouter(repo *database.UserEventsRepository) *mux.Router {
    r := mux.NewRouter()

    r.HandleFunc("/api/v1/user-action", middleware.UserEventsValidationMiddleware(api.HandleUserAction(repo))).Methods("POST")
    r.HandleFunc("/api/v1/user-actions", api.GetAllUserActions(repo)).Methods("GET")
    r.HandleFunc("/api/v1/users-actions/", api.GetUserActions(repo)).Methods("GET")
    r.HandleFunc("/api/v1/action-and-device-counts", api.GetActionAndDeviceCounts(repo)).Methods("GET")

    return r
}
