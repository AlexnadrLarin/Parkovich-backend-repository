package middleware

import (
    "encoding/json"
    "net/http"
)

type UserAction struct {
    EventType  string `json:"event_type"`
    UserAgent  string `json:"user_agent"`
    DeviceType string `json:"device_type"`
    EventTime  string `json:"event_time"` 
}

func UserEventsValidationMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var action UserAction
        
        if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
            http.Error(w, "Неверные данные", http.StatusBadRequest)
            return
        }

        if action.EventType == "" || action.UserAgent == "" || action.DeviceType == "" {
            http.Error(w, "Отсутствуют необходимые поля", http.StatusBadRequest)
            return
        }

        validEventTypes := []string{"session_scrolled_1", "session_scrolled_2", "session_scrolled_3", 
                                    "session_scrolled_4", "session_scrolled_5", "session_scrolled_6",
                                    "session_scrolled_7", "session_scrolled_8", "visited" ,"click_try", "comment"}
        isValid := false
        for _, validType := range validEventTypes {
            if action.EventType == validType {
                isValid = true
                break
            }
        }
        if !isValid {
            http.Error(w, "Недопустимый тип действия", http.StatusBadRequest)
            return
        }

        if action.UserAgent == "" {
            http.Error(w, "Не указана информация о браузере", http.StatusBadRequest)
            return
        }

        if action.DeviceType != "desktop" && action.DeviceType != "mobile" {
            http.Error(w, "Тип устройства должен быть 'desktop' или 'mobile'", http.StatusBadRequest)
            return
        }

        next.ServeHTTP(w, r)
    }
}
