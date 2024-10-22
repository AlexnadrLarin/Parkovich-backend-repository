package database

import (
    "context"
    "log"
    "os"
    "time"
    "fmt"

    "github.com/jackc/pgx/v4"
    "github.com/joho/godotenv"
)

type UserMessage struct {
    ID        int
    Name      string
    Email     string
    Message   string
    CreatedAt time.Time
}

type EmailSubscriber struct {
    ID           int
    Email        string
    SubscribedAt time.Time
}

type UserMessagesRepository struct {
    db *pgx.Conn
}

func NewUserMessagesRepository() (*UserMessagesRepository, error) {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Ошибка при загрузке файла .env: %v", err)
    }

    conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        return nil, err
    }

    return &UserMessagesRepository{db: conn}, nil
}

func (repo *UserMessagesRepository) SaveUserMessage(name, email, message string) error {
    tx, err := repo.db.Begin(context.Background())
    if err != nil {
        return err
    }

    defer func() {
        if err != nil {
            _ = tx.Rollback(context.Background())
        } else {
            _ = tx.Commit(context.Background())
        }
    }()

    queryInsertMessage := `INSERT INTO user_messages (name, email, message, created_at) VALUES ($1, $2, $3, $4)`
    _, err = tx.Exec(context.Background(), queryInsertMessage, name, email, message, time.Now())
    if err != nil {
        return err
    }

    queryCheck := `SELECT COUNT(*) FROM email_subscribers WHERE email = $1`
    var count int
    err = tx.QueryRow(context.Background(), queryCheck, email).Scan(&count)
    if err != nil {
        log.Printf("Ошибка при проверке существующего email: %v", err)
        return err
    }

    if count > 0 {
        return nil
    }

    queryInsertSubscriber := `INSERT INTO email_subscribers (email, subscribed_at) 
                              VALUES ($1, $2) 
                              ON CONFLICT (email) DO NOTHING`
    _, err = tx.Exec(context.Background(), queryInsertSubscriber, email, time.Now())
    if err != nil {
        return err
    }

    return nil
}

func (repo *UserMessagesRepository) SaveEmailSubscriber(email string) error {
    tx, err := repo.db.Begin(context.Background())
    if err != nil {
        return err
    }
    
    defer func() {
        if err != nil {
            _ = tx.Rollback(context.Background())
        } else {
            _ = tx.Commit(context.Background())
        }
    }()
    
    queryCheck := `SELECT COUNT(*) FROM email_subscribers WHERE email = $1`
    var count int
    err = tx.QueryRow(context.Background(), queryCheck, email).Scan(&count)
    if err != nil {
        log.Printf("Ошибка при проверке существующего email: %v", err)
        return err
    }

    if count > 0 {
        return fmt.Errorf("Подписчик уже есть в базе данных")
    }

    queryInsert := `INSERT INTO email_subscribers (email, subscribed_at) 
                    VALUES ($1, $2) 
                    ON CONFLICT (email) DO NOTHING`
    _, err = tx.Exec(context.Background(), queryInsert, email, time.Now())
    if err != nil {
        log.Printf("Ошибка при сохранении email в таблицу подписчиков: %v", err)
        return err
    }

    return nil
}


func (repo *UserMessagesRepository) GetUserMessages() ([]UserMessage, error) {
    query := `SELECT id, name, email, message, created_at FROM user_messages ORDER BY created_at DESC`
    rows, err := repo.db.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []UserMessage
    for rows.Next() {
        var msg UserMessage
        if err := rows.Scan(&msg.ID, &msg.Name, &msg.Email, &msg.Message, &msg.CreatedAt); err != nil {
            return nil, err
        }
        messages = append(messages, msg)
    }

    return messages, nil
}

func (repo *UserMessagesRepository) GetUserMessageByID(id int) (*UserMessage, error) {
    query := `SELECT id, name, email, message, created_at FROM user_messages WHERE id = $1`
    var message UserMessage
    err := repo.db.QueryRow(context.Background(), query, id).Scan(&message.ID, &message.Name, &message.Email, &message.Message, &message.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &message, nil
}

func (repo *UserMessagesRepository) GetEmailSubscribers() ([]EmailSubscriber, error) {
    query := `SELECT id, email, subscribed_at FROM email_subscribers ORDER BY subscribed_at DESC`
    rows, err := repo.db.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var subscribers []EmailSubscriber
    for rows.Next() {
        var sub EmailSubscriber
        if err := rows.Scan(&sub.ID, &sub.Email, &sub.SubscribedAt); err != nil {
            return nil, err
        }
        subscribers = append(subscribers, sub)
    }

    return subscribers, nil
}

func (repo *UserMessagesRepository) GetEmailSubscriberByID(id int) (*EmailSubscriber, error) {
    query := `SELECT id, email, subscribed_at FROM email_subscribers WHERE id = $1`
    var subscriber EmailSubscriber
    err := repo.db.QueryRow(context.Background(), query, id).Scan(&subscriber.ID, &subscriber.Email, &subscriber.SubscribedAt)
    if err != nil {
        return nil, err
    }
    return &subscriber, nil
}

func (repo *UserMessagesRepository) Close() {
    repo.db.Close(context.Background())
}
