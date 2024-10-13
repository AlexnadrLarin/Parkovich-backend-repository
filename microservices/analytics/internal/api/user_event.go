package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go-parkovich/microservices/analytics/internal/database"
	"go-parkovich/microservices/analytics/pkg/proto"
)

type UserEventServer struct {
    userevents.UnimplementedUserEventServiceServer
    repo       *database.UserEventsRepository  
    grpcClient userevents.UserEventServiceClient   
}

func NewUserEventServer(repo *database.UserEventsRepository, grpcClient userevents.UserEventServiceClient) *UserEventServer {
    return &UserEventServer{
        repo: repo,
        grpcClient: grpcClient,
    }
}

type UserAction struct {
    UserID     string    `json:"user_id"`
    EventType  string    `json:"event_type"`
    UserAgent  string    `json:"user_agent"`
    DeviceType string    `json:"device_type"`
    EventTime  time.Time `json:"event_time"`
}

func HandleUserAction(repo *database.UserEventsRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var action UserAction
        if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
            respondWithJSON(w, http.StatusBadRequest, "Неверные данные")
            return
        }

        if action.EventTime.IsZero() {
            action.EventTime = time.Now()
        }

        dbAction := &database.UserAction{
            EventType:  action.EventType,
            UserAgent:  action.UserAgent,
            DeviceType: action.DeviceType,
            EventTime:  action.EventTime,
        }

        if err := repo.SaveOrUpdateUserAction(dbAction); err != nil {
            log.Printf("Ошибка при сохранении действия: %v", err)
            respondWithJSON(w, http.StatusInternalServerError, "Ошибка при сохранении данных")
            return
        }

        respondWithJSON(w, http.StatusOK, "Действие сохранено")
    }
}

func GetAllUserActions(repo *database.UserEventsRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        actions, err := repo.GetAllUserActions()
        if err != nil {
            log.Printf("Ошибка при получении действий: %v", err)
            respondWithJSON(w, http.StatusInternalServerError, "Ошибка при получении данных")
            return
        }

        respondWithJSON(w, http.StatusOK, actions)
    }
}

func GetUserActions(repo *database.UserEventsRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userIDStr := r.URL.Query().Get("user_id")
        if userIDStr == "" {
            respondWithJSON(w, http.StatusBadRequest, "Необходимо указать user_id")
            return
        }

        actions, err := repo.GetUserActionsByID(userIDStr)
        if err != nil {
            log.Printf("Ошибка при получении действий: %v", err)
            respondWithJSON(w, http.StatusInternalServerError, "Ошибка при получении данных")
            return
        }

        respondWithJSON(w, http.StatusOK, actions)
    }
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}
