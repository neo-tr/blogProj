package models

import (
	"blogProj/db"
	"context"
	"time"
)

type User struct {
	ID        int
	Username  string
	Password  string // хранится хэш
	Nickname  string
	CreatedAt time.Time
}

// CreateUser создаёт нового пользователя
func CreateUser(username, hashedPassword, nickname string) error {
	sql := INSERT INTO users (username, password, nickname) VALUES ($1, $2, $3)
	_, err := db.Pool.Exec(context.Background(), sql, username, hashedPassword, nickname)
	return err
}

// GetUserByUsername возвращает пользователя по логину
func GetUserByUsername(username string) (User, error) {
	sql := SELECT id, username, password, nickname, created_at FROM users WHERE username=$1
	var u User
	err := db.Pool.QueryRow(context.Background(), sql, username).Scan(&u.ID, &u.Username, &u.Password, &u.Nickname, &u.CreatedAt)
	return u, err
}