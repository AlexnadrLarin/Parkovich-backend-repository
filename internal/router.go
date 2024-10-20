package internal

import (
    "net/http"
	
	"github.com/swaggo/http-swagger"
    "github.com/gorilla/mux"

    "go-parkovich/internal/middleware"
    "go-parkovich/internal/api"
    "go-parkovich/internal/database"

)

func SetupRouter(repo *database.UserMessagesRepository) *mux.Router {
    r := mux.NewRouter()

    r.PathPrefix("/swagger/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        httpSwagger.WrapHandler.ServeHTTP(w, r)
    })

    r.HandleFunc("/api/v1/user-message", middleware.UserMessageValidation(api.SaveUserMessage(repo))).Methods("POST")
    r.HandleFunc("/api/v1/email-subscribe", middleware.EmailSubscriberValidation(api.SaveEmailSubscriber(repo))).Methods("POST")
    r.HandleFunc("/api/v1/user-messages", api.GetAllUserMessages(repo)).Methods("GET")
    r.HandleFunc("/api/v1/user-messages/{id}", api.GetUserMessageByID(repo)).Methods("GET")
    r.HandleFunc("/api/v1/email-subscribers", api.GetAllEmailSubscribers(repo)).Methods("GET")
    r.HandleFunc("/api/v1/email-subscribers/{id}", api.GetEmailSubscriberByID(repo)).Methods("GET")

    return r
}
