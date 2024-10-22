package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
    "io"
    "bytes"
    "fmt"

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

// UserAction представляет действие пользователя
type UserAction struct {
    UserID     string    `json:"user_id"`
    EventType  string    `json:"event_type"`
    UserAgent  string    `json:"user_agent"`
    DeviceType string    `json:"device_type"`
    EventTime  time.Time `json:"event_time"`
}

// HandleUserAction сохраняет действие пользователя
// @Summary Сохранение действия пользователя
// @Description Сохраняет действие пользователя в базе данных
// @Tags UserAction
// @Accept json
// @Produce json
// @Param action body UserAction true "Данные действия пользователя"
// @Success 200 {string} string "Действие сохранено"
// @Failure 400 {string} string "Неверные данные"
// @Failure 500 {string} string "Ошибка при сохранении данных"
// @Router /api/v1/user-action [post]
func HandleUserAction(repo *database.UserEventsRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var action UserAction
        
        body, err := io.ReadAll(r.Body)
        if err != nil {
            respondWithJSON(w, http.StatusInternalServerError, ErrorMessage("Ошибка при чтении данных"))
            return
        }

        r.Body = io.NopCloser(bytes.NewBuffer(body))

        if err := json.Unmarshal(body, &action);  err != nil {
            respondWithJSON(w, http.StatusBadRequest, ErrorMessage("Неверные данные"))
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
            respondWithJSON(w, http.StatusInternalServerError, ErrorMessage("Ошибка при сохранении данных"))
            return
        }
        
        log.Println("Действие сохранено")
        respondWithJSON(w, http.StatusOK, SuccessMessage("Действие сохранено"))
    }
}

// GetAllUserActions Предоставляет данные о действиях всех пользователей
// @Summary Предоставление информации о действиях всех пользователей
// @Description Предоставление информации о действиях всех пользователей из базы данных
// @Tags UserAction
// @Produce json
// @Param action body UserAction true "Данные действия пользователя"
// @Success 200 {string} json Предоставляются данные о действиях пользователя из таблицы
// @Failure 500 {string} string "Ошибка при получении действий"
// @Router /api/v1/user-actions [get]
func GetAllUserActions(repo *database.UserEventsRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        actions, err := repo.GetAllUserActions()
        if err != nil {
            log.Printf("Ошибка при получении действий: %v", err)
            respondWithJSON(w, http.StatusInternalServerError, ErrorMessage("Ошибка при получении данных"))
            return
        }

        respondWithJSON(w, http.StatusOK, actions)
    }
}

// GetUserByID возвращает данные о конкретном пользователе по его ID
// @Summary Получить данные о пользователе
// @Description Возвращает данные о конкретном пользователе по его уникальному идентификатору (user_id)
// @Tags UserAction
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Success 200 {string} json "Данные о пользователе"
// @Failure 400 {string} string "Некорректный запрос, не указан user_id"
// @Failure 404 {string} string "Пользователь с таким ID не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/v1/users-actions/{user_id} [get]
func GetUserActions(repo *database.UserEventsRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userIDStr := r.URL.Query().Get("user_id")
        if userIDStr == "" {
            respondWithJSON(w, http.StatusBadRequest, ErrorMessage("Необходимо указать user_id"))
            return
        }

        actions, err := repo.GetUserActionsByID(userIDStr)
        if err != nil {
            log.Printf("Ошибка при получении действий: %v", err)
            respondWithJSON(w, http.StatusInternalServerError, ErrorMessage("Ошибка при получении данных"))
            return
        }

        if len(actions) == 0 {
            log.Printf("Пользователь с таким ID не найден")
            respondWithJSON(w, http.StatusNotFound, ErrorMessage("Пользователь не найден"))
            return
        }

        respondWithJSON(w, http.StatusOK, actions)
    }
}

// GetActionAndDeviceCountsHandler возвращает количество действий по типам, браузерам и устройствам
// @Summary Получить количество действий по типам, браузерам и устройствам
// @Description Возвращает количество действий, сгруппированных по типу действия, браузеру и типу устройства
// @Tags User Actions
// @Produce json
// @Success 200 {object} map[string]interface{} "Количество действий по типам, браузерам и устройствам"
// @Failure 500 {string} string "Ошибка при получении данных"
// @Router /api/v1/action-and-device-counts [get]
func GetActionAndDeviceCounts(repo *database.UserEventsRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        actionCounts, userAgentCounts, deviceTypeCounts, err := repo.GetActionCountsByType()
        if err != nil {
            log.Printf("Ошибка при подсчете действий, браузеров и устройств: %v", err)
            respondWithJSON(w, http.StatusInternalServerError, ErrorMessage("Ошибка при получении данных"))
            return
        }

        visitors := 0
        actionTry := 0
        actionMessage := 0
        sections := make([]int, 8) 

        for eventType, count := range actionCounts {
            switch eventType {
            case "visited":
                visitors = int(count)
            case "click_try":
                actionTry = int(count)
            case "comment":
                actionMessage = int(count)
            default:
                for i := 1; i <= 8; i++ {
                    if eventType == fmt.Sprintf("session_scrolled_%d", i) {
                        sections[i-1] = int(count)
                    }
                }
            }
        }

        result := map[string]interface{}{
            "visitors":       visitors,
            "devices":        deviceTypeCounts,  
            "browsers":       userAgentCounts,  
            "action_try":     actionTry,
            "action_message": actionMessage,
            "sections":       sections,
        }

        respondWithJSON(w, http.StatusOK, result)
    }
}

func ErrorMessage(message string) map[string]string {
    return map[string]string{"error": message}
}

func SuccessMessage(message string) map[string]string {
    return map[string]string{"message": message}
}

func respondWithJSON(w http.ResponseWriter, statusCode int,  data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)

    json.NewEncoder(w).Encode(data)
}
