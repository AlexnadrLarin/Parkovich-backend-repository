package api

import (
    "encoding/json"
    "log"
    "net/http"
	"strconv"
	"io"
	"bytes"

    "go-parkovich/internal/middleware"
    "go-parkovich/internal/database"
)

type UserMessage struct {
    Name    string `json:"name"`
    Email   string `json:"email"`
    Message string `json:"message"`
}

type EmailSubscriber struct {
    Email string `json:"email"`
}

// SaveUserMessage сохраняет сообщение пользователя
// @Summary Сохранение сообщения пользователя
// @Description Сохраняет сообщение пользователя с валидацией имени, email и сообщения
// @Tags User Messages
// @Accept json
// @Produce json
// @Param message body UserMessage true "Данные сообшения пользователя"
// @Success 200 {string} string "Сообщение успешно сохранено"
// @Failure 400 {string} string "Некорректные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/v1/user-message [post]
func SaveUserMessage(repo *database.UserMessagesRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var message middleware.UserMessage 
		
		body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Ошибка при чтении данных", http.StatusInternalServerError)
            return
        }

        r.Body = io.NopCloser(bytes.NewBuffer(body))


        if err := json.Unmarshal(body, &message); err != nil {
			log.Printf("%v", err)
            http.Error(w, "Неверные данные", http.StatusBadRequest)
            return
        }

        err = repo.SaveUserMessage(message.Name, message.Email, message.Message)
        if err != nil {
            log.Printf("Ошибка при сохранении сообщения: %v", err)
            http.Error(w, "Ошибка при сохранении сообщения", http.StatusInternalServerError)
            return
        }

        err = repo.SaveEmailSubscriber(message.Email)
        if err != nil {
            log.Printf("Ошибка при сохранении подписчика: %v", err)
            http.Error(w, "Ошибка при сохранении подписчика", http.StatusInternalServerError)
            return
        }

        log.Println("Сообщение и подписка успешно сохранены")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Сообщение и подписка успешно сохранены"))
    }
}

// SaveEmailSubscriber сохраняет email подписчика
// @Summary Сохранение email подписчика
// @Description Сохраняет email подписчика с валидацией
// @Tags Email Subscribers
// @Accept json
// @Produce json
// @Param subcriber body EmailSubscriber true "Данные email пользователя"
// @Success 200 {string} string "Подписчик успешно сохранён"
// @Failure 400 {string} string "Некорректные данные"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/v1/email-subscribe [post]
func SaveEmailSubscriber(repo *database.UserMessagesRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var subscriber middleware.EmailSubscriber

		body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Ошибка при чтении данных", http.StatusInternalServerError)
            return
        }

        r.Body = io.NopCloser(bytes.NewBuffer(body))

        if err := json.Unmarshal(body, &subscriber); err != nil {
            http.Error(w, "Неверные данные", http.StatusBadRequest)
            return
        }

        err = repo.SaveEmailSubscriber(subscriber.Email)
        if err != nil {
            log.Printf("Ошибка при сохранении подписчика: %v", err)
            http.Error(w, "Ошибка при сохранении подписчика", http.StatusInternalServerError)
            return
        }

        log.Println("Подписчик успешно сохранен")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Подписчик успешно сохранен"))
    }
}

// GetAllUserMessages возвращает все сообщения пользователей
// @Summary Получить все сообщения пользователей
// @Description Возвращает все сообщения из таблицы сообщений пользователей
// @Tags User Messages
// @Produce json
// @Success 200 {object} array "Все сообщения пользователей"
// @Failure 500 {string} string "Ошибка при получении сообщений"
// @Router /api/v1/user-messages [get]
func GetAllUserMessages(repo *database.UserMessagesRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        messages, err := repo.GetUserMessages()
        if err != nil {
            log.Printf("Ошибка при получении сообщений пользователей: %v", err)
            http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
            return
        }
        respondWithJSON(w, http.StatusOK, messages)
    }
}

// GetUserMessageByID возвращает сообщение пользователя по ID
// @Summary Получить сообщение пользователя по ID
// @Description Возвращает конкретное сообщение пользователя по его ID
// @Tags User Messages
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string "Сообщение пользователя"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 404 {string} string "Сообщение не найдено"
// @Router /api/v1/user-messages/{id} [get]
func GetUserMessageByID(repo *database.UserMessagesRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := r.URL.Query().Get("id")
        id, err := strconv.Atoi(idStr)
        if err != nil || id <= 0 {
            http.Error(w, "Некорректный ID", http.StatusBadRequest)
            return
        }

        message, err := repo.GetUserMessageByID(id)
        if err != nil {
            log.Printf("Сообщение не найдено: %v", err)
            http.Error(w, "Сообщение не найдено", http.StatusNotFound)
            return
        }
        respondWithJSON(w, http.StatusOK, message)
    }
}

// GetAllEmailSubscribers возвращает всех подписчиков
// @Summary Получить всех подписчиков
// @Description Возвращает всех подписчиков из таблицы email_subscribers
// @Tags Email Subscribers
// @Produce json
// @Success 200 {object} array "Все подписчики"
// @Failure 500 {string} string "Ошибка при получении подписчиков"
// @Router /api/v1/email-subscribers [get]
func GetAllEmailSubscribers(repo *database.UserMessagesRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        subscribers, err := repo.GetEmailSubscribers()
        if err != nil {
            log.Printf("Ошибка при получении подписчиков: %v", err)
            http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
            return
        }
        respondWithJSON(w, http.StatusOK, subscribers)
    }
}

// GetEmailSubscriberByID возвращает подписчика по ID
// @Summary Получить подписчика по ID
// @Description Возвращает конкретного подписчика по его ID
// @Tags Email Subscribers
// @Produce json
// @Param id path string true "ID"
// @Success 200 {string} string "Подписчик"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 404 {string} string "Подписчик не найден"
// @Router /api/v1/email-subscribers/{id} [get]
func GetEmailSubscriberByID(repo *database.UserMessagesRepository) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := r.URL.Query().Get("id")
        id, err := strconv.Atoi(idStr)
        if err != nil || id <= 0 {
            http.Error(w, "Некорректный ID", http.StatusBadRequest)
            return
        }

        subscriber, err := repo.GetEmailSubscriberByID(id)
        if err != nil {
            log.Printf("Подписчик не найден: %v", err)
            http.Error(w, "Подписчик не найден", http.StatusNotFound)
            return
        }
        respondWithJSON(w, http.StatusOK, subscriber)
    }
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}