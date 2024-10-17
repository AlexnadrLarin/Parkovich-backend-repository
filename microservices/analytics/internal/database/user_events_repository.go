package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
    "github.com/joho/godotenv"
)

type UserAction struct {
    UserID     uint64
    EventType  string
    UserAgent  string
    DeviceType string
    EventTime  time.Time
}

type UserEventsRepository struct {
    db clickhouse.Conn
}

func NewUserEventsRepository() (*UserEventsRepository, error) {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Ошибка при загрузке файла .env: %v", err)
    }

    conn, err := clickhouse.Open(&clickhouse.Options{
        Addr: []string{os.Getenv("CLICKHOUSE_HOST") + ":" + os.Getenv("CLICKHOUSE_PORT")},
        Auth: clickhouse.Auth{
            Username: os.Getenv("CLICKHOUSE_USERNAME"),
            Password: os.Getenv("CLICKHOUSE_PASSWORD"),
            Database: os.Getenv("CLICKHOUSE_DATABASE"),
        },
    })
    if err != nil {
        return nil, err
    }

    return &UserEventsRepository{db: conn}, nil
}

func (repo *UserEventsRepository) SaveOrUpdateUserAction(action *UserAction) error {
    queryCheck := `SELECT COUNT(*) FROM user_events WHERE user_id = ? AND event_type = ?`
    
    var count uint64
    err := repo.db.QueryRow(context.Background(), queryCheck, action.UserID, action.EventType).Scan(&count)
    if err != nil {
        log.Printf("Ошибка при проверке существования действия: %v", err)
        return err
    }

    if count > 0 {
        queryUpdate := `UPDATE user_events SET action_count = action_count + 1 WHERE user_id = ? AND event_type = ?`
        err := repo.db.Exec(context.Background(), queryUpdate, action.UserID, action.EventType)
        if err != nil {
            log.Printf("Ошибка при обновлении счетчика действий: %v", err)
            return err
        }
    } else {
        queryInsert := `INSERT INTO user_events (user_id, event_type, user_agent, device_type, event_time, action_count) 
                        VALUES (?, ?, ?, ?, ?, 1)`

        var lastUserID uint64
        err := repo.db.QueryRow(context.Background(), "SELECT max(user_id) FROM user_events").Scan(&lastUserID)
        if err != nil {
            log.Printf("Ошибка при получении последнего user_id: %v", err)
            return err
        }
        newUserID := lastUserID + 1

        err = repo.db.Exec(context.Background(), queryInsert, newUserID, action.EventType, action.UserAgent, action.DeviceType, action.EventTime)
        if err != nil {
            log.Printf("Ошибка вставки данных: %v", err)
            return err
        }
    }
    return nil
}

func (repo *UserEventsRepository) GetAllUserActions() ([]UserAction, error) {
    query := `SELECT user_id, event_type, user_agent, device_type, event_time FROM user_events`
    rows, err := repo.db.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var actions []UserAction
    for rows.Next() {
        var action UserAction
        if err := rows.Scan(&action.UserID, &action.EventType, &action.UserAgent, &action.DeviceType, &action.EventTime); err != nil {
            return nil, err
        }
        actions = append(actions, action)
    }
    return actions, nil
}

func (repo *UserEventsRepository) GetUserActionsByID(userID string) ([]UserAction, error) {
    query := `SELECT user_id, event_type, user_agent, device_type, event_time FROM user_events WHERE user_id = ?`
    rows, err := repo.db.Query(context.Background(), query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var actions []UserAction
    for rows.Next() {
        var action UserAction
        if err := rows.Scan(&action.UserID, &action.EventType, &action.UserAgent, &action.DeviceType, &action.EventTime); err != nil {
            return nil, err
        }
        actions = append(actions, action)
    }
    return actions, nil
}

func (repo *UserEventsRepository) GetActionCountsByType() (map[string]uint64, error) {
    query := `SELECT event_type, COUNT(*) FROM user_events GROUP BY event_type`
    rows, err := repo.db.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    actionCounts := make(map[string]uint64)
    for rows.Next() {
        var eventType string
        var count uint64
        if err := rows.Scan(&eventType, &count); err != nil {
            return nil, err
        }
        actionCounts[eventType] = count
    }

    return actionCounts, nil
}

func (repo *UserEventsRepository) Close() {
    repo.db.Close()
}