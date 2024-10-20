package middleware

import (
    "log"
    "net/http"
    "regexp"
    "encoding/json"
	"unicode"
	"io"
	"bytes"
)

type UserMessage struct {
    Name    string `json:"name"`
    Email   string `json:"email"`
    Message string `json:"message"`
}

type EmailSubscriber struct {
    Email string `json:"email"`
}

func ValidateEmail(email string) bool {
    var re = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

func ContainsZalgoText(input string) bool {
    for _, r := range input {
        if r >= '\u0300' && r <= '\u036F' { 
            return true
        }
        if unicode.IsControl(r) { 
            return true
        }
    }
    return false
}

func ValidateName(name string) bool {
    re := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9\s]+$`)
    return re.MatchString(name) && !ContainsZalgoText(name)
}

func ValidateMessage(message string) bool {
    return len(message) > 0 && len(message) <= 500 && !ContainsZalgoText(message)
}

func UserMessageValidation(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var message UserMessage

		body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Ошибка при чтении данных", http.StatusInternalServerError)
            return
        }

        r.Body = io.NopCloser(bytes.NewBuffer(body))

		if err := json.Unmarshal(body, &message); err != nil {
            http.Error(w, "Неверные данные", http.StatusBadRequest)
            return
        }

        if !ValidateName(message.Name) {
            log.Println("Некорректное имя:", message.Name)
            http.Error(w, "Некорректное имя", http.StatusBadRequest)
            return
        }

        if !ValidateMessage(message.Message) {
            log.Println("Некорректное сообщение:", message.Message)
            http.Error(w, "Сообщение должно быть не пустым, длиной до 500 символов и без Zalgo-текста", http.StatusBadRequest)
            return
        }

        if !ValidateEmail(message.Email) {
            log.Println("Некорректный email:", message.Email)
            http.Error(w, "Некорректный email", http.StatusBadRequest)
            return
        }

        if !ValidateMessage(message.Message) {
            log.Println("Некорректное сообщение:", message.Message)
            http.Error(w, "Сообщение должно быть не пустым и длиной не более 500 символов", http.StatusBadRequest)
            return
        }

        next.ServeHTTP(w, r)
    }
}

func EmailSubscriberValidation(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var subscriber EmailSubscriber

		body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Ошибка при чтении данных", http.StatusInternalServerError)
            return
        }

        r.Body = io.NopCloser(bytes.NewBuffer(body))
		
        if err := json.Unmarshal(body, &subscriber);  err != nil {
            http.Error(w, "Неверные данные", http.StatusBadRequest)
            return
        }

        if !ValidateEmail(subscriber.Email) {
            log.Println("Некорректный email:", subscriber.Email)
            http.Error(w, "Некорректный email", http.StatusBadRequest)
            return
        }

        next.ServeHTTP(w, r)
    }
}
